package gateway

import (
	"bytes"
	"context"
	"net"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
	"text/template"

	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/config"
	"github.com/TykTechnologies/tyk/storage"
)

const sampleUptimeTestAPI = `{
	"api_id": "test",
	"use_keyless": true,
	"version_data": {
		"not_versioned": true,
		"versions": {
			"v1": {"name": "v1"}
		}
	},
	"uptime_tests": {
		"check_list": [
			{
				"url": "{{.Host1}}/get",
				"method": "GET"
			},
			{
				"url": "{{.Host2}}/get",
				"method": "GET"
			}
		]
	},
	"proxy": {
		"listen_path": "/",
		"enable_load_balancing": true,
		"check_host_against_uptime_tests": true,
		"target_list": [
			"{{.Host1}}",
			"{{.Host2}}"
		]
	}
}`

type testEventHandler struct {
	cb func(config.EventMessage)
}

func (w *testEventHandler) Init(handlerConf interface{}) error {
	return nil
}

func (w *testEventHandler) HandleEvent(em config.EventMessage) {
	w.cb(em)
}

func TestHostChecker(t *testing.T) {
	specTmpl := template.Must(template.New("spec").Parse(sampleUptimeTestAPI))

	tmplData := struct {
		Host1, Host2 string
	}{
		testHttpAny,
		testHttpFailureAny,
	}

	specBuf := &bytes.Buffer{}
	specTmpl.ExecuteTemplate(specBuf, specTmpl.Name(), &tmplData)

	spec := CreateDefinitionFromString(specBuf.String())

	// From api_loader.go#processSpec
	sl := apidef.NewHostListFromList(spec.Proxy.Targets)
	spec.Proxy.StructuredTargetList = sl

	var eventWG sync.WaitGroup
	// Should receive one HostDown event
	eventWG.Add(1)
	cb := func(em config.EventMessage) {
		eventWG.Done()
	}

	spec.EventPaths = map[apidef.TykEvent][]config.TykEventHandler{
		"HostDown": {&testEventHandler{cb}},
	}

	apisMu.Lock()
	apisByID = map[string]*APISpec{spec.APIID: spec}
	apisMu.Unlock()
	GlobalHostChecker.checkerMu.Lock()
	GlobalHostChecker.checker.sampleTriggerLimit = 1
	GlobalHostChecker.checkerMu.Unlock()
	defer func() {
		apisMu.Lock()
		apisByID = make(map[string]*APISpec)
		apisMu.Unlock()
		GlobalHostChecker.checkerMu.Lock()
		GlobalHostChecker.checker.sampleTriggerLimit = defaultSampletTriggerLimit
		GlobalHostChecker.checkerMu.Unlock()
	}()

	SetCheckerHostList()
	GlobalHostChecker.checkerMu.Lock()
	if len(GlobalHostChecker.currentHostList) != 2 {
		t.Error("Should update hosts manager check list", GlobalHostChecker.currentHostList)
	}

	if len(GlobalHostChecker.checker.newList) != 2 {
		t.Error("Should update host checker check list")
	}
	GlobalHostChecker.checkerMu.Unlock()

	hostCheckTicker <- struct{}{}
	eventWG.Wait()

	if GlobalHostChecker.HostDown(testHttpAny) {
		t.Error("Should not mark as down")
	}

	if !GlobalHostChecker.HostDown(testHttpFailureAny) {
		t.Error("Should mark as down")
	}

	// Test it many times concurrently, to simulate concurrent and
	// parallel requests to the API. This will catch bugs in those
	// scenarios, like data races.
	var targetWG sync.WaitGroup
	for i := 0; i < 10; i++ {
		targetWG.Add(1)
		go func() {
			host, err := nextTarget(spec.Proxy.StructuredTargetList, spec)
			if err != nil {
				t.Error("Should return nil error, got", err)
			}
			if host != testHttpAny {
				t.Error("Should return only active host, got", host)
			}
			targetWG.Done()
		}()
	}
	targetWG.Wait()

	GlobalHostChecker.checkerMu.Lock()
	if GlobalHostChecker.checker.checkTimeout != defaultTimeout {
		t.Error("Should set defaults", GlobalHostChecker.checker.checkTimeout)
	}

	redisStore := GlobalHostChecker.store.(*storage.RedisCluster)
	if ttl, _ := redisStore.GetKeyTTL(PoolerHostSentinelKeyPrefix + testHttpFailure); int(ttl) != GlobalHostChecker.checker.checkTimeout*GlobalHostChecker.checker.sampleTriggerLimit {
		t.Error("HostDown expiration key should be checkTimeout + 1", ttl)
	}
	GlobalHostChecker.checkerMu.Unlock()
}

func TestReverseProxyAllDown(t *testing.T) {
	specTmpl := template.Must(template.New("spec").Parse(sampleUptimeTestAPI))

	tmplData := struct {
		Host1, Host2 string
	}{
		testHttpFailureAny,
		testHttpFailureAny,
	}

	specBuf := &bytes.Buffer{}
	specTmpl.ExecuteTemplate(specBuf, specTmpl.Name(), &tmplData)

	spec := CreateDefinitionFromString(specBuf.String())

	// From api_loader.go#processSpec
	sl := apidef.NewHostListFromList(spec.Proxy.Targets)
	spec.Proxy.StructuredTargetList = sl

	var eventWG sync.WaitGroup
	// Should receive one HostDown event
	eventWG.Add(1)
	cb := func(em config.EventMessage) {
		eventWG.Done()
	}
	spec.EventPaths = map[apidef.TykEvent][]config.TykEventHandler{
		"HostDown": {&testEventHandler{cb}},
	}

	apisMu.Lock()
	apisByID = map[string]*APISpec{spec.APIID: spec}
	apisMu.Unlock()
	GlobalHostChecker.checkerMu.Lock()
	GlobalHostChecker.checker.sampleTriggerLimit = 1
	GlobalHostChecker.checkerMu.Unlock()
	defer func() {
		apisMu.Lock()
		apisByID = make(map[string]*APISpec)
		apisMu.Unlock()
		GlobalHostChecker.checkerMu.Lock()
		GlobalHostChecker.checker.sampleTriggerLimit = defaultSampletTriggerLimit
		GlobalHostChecker.checkerMu.Unlock()
	}()

	SetCheckerHostList()

	hostCheckTicker <- struct{}{}
	eventWG.Wait()

	remote, _ := url.Parse(testHttpAny)
	proxy := TykNewSingleHostReverseProxy(remote, spec)

	req := TestReq(t, "GET", "/", nil)
	rec := httptest.NewRecorder()
	proxy.ServeHTTP(rec, req)
	if rec.Code != 503 {
		t.Fatalf("wanted code to be 503, was %d", rec.Code)
	}
}

func TestTestCheckerTCPHosts_correct_answers(t *testing.T) {
	log.Out = os.Stdout
	log.Level = 6
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	data := HostData{
		CheckURL: l.Addr().String(),
		Protocol: "tcp",
		Commands: []apidef.CheckCommand{
			{
				Name: "send", Message: "ping",
			}, {
				Name: "expect", Message: "pong",
			},
		},
	}
	go func(ls net.Listener) {
		for {
			s, err := ls.Accept()
			if err != nil {
				return
			}
			buf := make([]byte, 4)
			_, err = s.Read(buf)
			if err != nil {
				return
			}
			if string(buf) == "ping" {
				s.Write([]byte("pong"))
			} else {
				s.Write([]byte("unknown"))
			}
		}
	}(l)
	ctx, cancel := context.WithCancel(context.Background())
	hs := &HostUptimeChecker{}
	failed := false
	up := false
	ping := false
	hs.Init(2, 1, 0, map[string]HostData{
		l.Addr().String(): data,
	},
		func(HostHealthReport) {
			failed = true
			cancel()
		},
		func(HostHealthReport) {
			up = true
			cancel()
		},
		func(HostHealthReport) {
			ping = true
			cancel()
		},
	)
	hs.sampleTriggerLimit = 1
	go hs.Start()
	hostCheckTicker <- struct{}{}
	<-ctx.Done()
	hs.Stop()
	if !(ping && !failed && !up) {
		t.Errorf("expected the host to be up : field:%v up:%v pinged:%v", failed, up, ping)
	}
}
func TestTestCheckerTCPHosts_correct_wrong_answers(t *testing.T) {
	log.Out = os.Stdout
	log.Level = 6
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	data := HostData{
		CheckURL: l.Addr().String(),
		Protocol: "tcp",
		Commands: []apidef.CheckCommand{
			{
				Name: "send", Message: "ping",
			}, {
				Name: "expect", Message: "pong",
			},
		},
	}
	go func(ls net.Listener) {
		for {
			s, err := ls.Accept()
			if err != nil {
				return
			}
			buf := make([]byte, 4)
			_, err = s.Read(buf)
			if err != nil {
				return
			}
			s.Write([]byte("unknown"))
		}
	}(l)
	ctx, cancel := context.WithCancel(context.Background())
	hs := &HostUptimeChecker{}
	failed := false
	up := false
	ping := false
	hs.Init(2, 1, 0, map[string]HostData{
		l.Addr().String(): data,
	},
		func(HostHealthReport) {
			failed = true
			cancel()
		},
		func(HostHealthReport) {
			up = true
			cancel()
		},
		func(HostHealthReport) {
			ping = true
			cancel()
		},
	)
	hs.sampleTriggerLimit = 1
	go hs.Start()
	hostCheckTicker <- struct{}{}
	<-ctx.Done()
	hs.Stop()
	if !failed {
		t.Errorf("expected the host to be down : field:%v up:%v pinged:%v", failed, up, ping)
	}
}
