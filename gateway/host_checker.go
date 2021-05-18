package gateway

import (
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/config"
	"github.com/jeffail/tunny"
	proxyproto "github.com/pires/go-proxyproto"
)

const (
	defaultTimeout             = 10
	defaultSampletTriggerLimit = 3
)

var (
	HostCheckerClient = &http.Client{
		Timeout: 500 * time.Millisecond,
	}
	defaultWorkerPoolSize = runtime.NumCPU()
	hostCheckTicker       = make(chan struct{})
)

type HostData struct {
	CheckURL            string
	Protocol            string
	Timeout             time.Duration
	EnableProxyProtocol bool
	Commands            []apidef.CheckCommand
	Method              string
	Headers             map[string]string
	Body                string
	MetaData            map[string]string
}

type HostHealthReport struct {
	HostData
	ResponseCode int
	Latency      float64
	IsTCPError   bool
}

type HostSample struct {
	count        int
	reachedLimit bool
}

type HostUptimeChecker struct {
	failureCallback    func(HostHealthReport)
	upCallback         func(HostHealthReport)
	pingCallback       func(HostHealthReport)
	workerPoolSize     int
	sampleTriggerLimit int
	checkTimeout       int
	HostList           map[string]HostData
	pool               *tunny.WorkPool

	errorChan       chan HostHealthReport
	okChan          chan HostHealthReport
	stopPollingChan chan bool
	samples         *sync.Map
	stopLoop        bool
	muStopLoop      sync.RWMutex

	resetListMu sync.Mutex
	doResetList bool
	newList     map[string]HostData
}

func (h *HostUptimeChecker) getStopLoop() bool {
	h.muStopLoop.RLock()
	defer h.muStopLoop.RUnlock()
	return h.stopLoop
}

func (h *HostUptimeChecker) setStopLoop(newValue bool) {
	h.muStopLoop.Lock()
	h.stopLoop = newValue
	h.muStopLoop.Unlock()
}

func (h *HostUptimeChecker) getStaggeredTime() time.Duration {
	if h.checkTimeout <= 5 {
		return time.Duration(h.checkTimeout) * time.Second
	}

	rand.Seed(time.Now().Unix())
	min := h.checkTimeout - 3
	max := h.checkTimeout + 3

	dur := rand.Intn(max-min) + min

	return time.Duration(dur) * time.Second
}

func (h *HostUptimeChecker) HostCheckLoop() {
	for !h.getStopLoop() {
		if isRunningTests() {
			<-hostCheckTicker
		}
		h.resetListMu.Lock()
		if h.doResetList && h.newList != nil {
			h.HostList = h.newList
			h.newList = nil
			h.doResetList = false
			log.Debug("[HOST CHECKER] Host list reset")
		}
		h.resetListMu.Unlock()
		for _, host := range h.HostList {
			_, err := h.pool.SendWork(host)
			if err != nil && err != tunny.ErrPoolNotRunning {
				log.Warnf("[HOST CHECKER] could not send work, error: %v", err)
			}
		}

		if !isRunningTests() {
			time.Sleep(h.getStaggeredTime())
		}
	}
	log.Info("[HOST CHECKER] Checker stopped")
}

func (h *HostUptimeChecker) HostReporter() {
	for {
		select {
		case okHost := <-h.okChan:
			// check if the the host url is in the sample map
			if hostSample, found := h.samples.Load(okHost.CheckURL); found {
				sample := hostSample.(HostSample)
				//if it reached the h.sampleTriggerLimit, we're going to start decreasing the count value
				if sample.reachedLimit {
					newCount := sample.count - 1

					if newCount <= 0 {
						//if the count-1 is equals to zero, it means that the host is fully up.

						h.samples.Delete(okHost.CheckURL)
						log.Warning("[HOST CHECKER] [HOST UP]: ", okHost.CheckURL)
						h.upCallback(okHost)
					} else {
						//in another case, we are one step closer. We just update the count number
						sample.count = newCount
						log.Warning("[HOST CHECKER] [HOST UP BUT NOT REACHED LIMIT]: ", okHost.CheckURL)
						h.samples.Store(okHost.CheckURL, sample)
					}
				}
			}
			go h.pingCallback(okHost)

		case failedHost := <-h.errorChan:
			sample := HostSample{
				count: 1,
			}

			//If a host fails, we check if it has failed already
			if hostSample, found := h.samples.Load(failedHost.CheckURL); found {
				sample = hostSample.(HostSample)
				// we add THIS failure to the count
				sample.count = sample.count + 1
			}

			if sample.count == h.sampleTriggerLimit {
				// if it reached the h.sampleTriggerLimit, it means the host is down for us. We update the reachedLimit flag and store it in the sample map
				log.Warning("[HOST CHECKER] [HOST DOWN]: ", failedHost.CheckURL)

				sample.reachedLimit = true
				h.samples.Store(failedHost.CheckURL, sample)
				go h.failureCallback(failedHost)

			} else if sample.count <= h.sampleTriggerLimit {
				//if it failed but not reached the h.sampleTriggerLimit yet, we just add the counter to the map.
				log.Warning("[HOST CHECKER] [HOST DOWN BUT NOT REACHED LIMIT]: ", failedHost.CheckURL)
				h.samples.Store(failedHost.CheckURL, sample)
			}

			go h.pingCallback(failedHost)

		case <-h.stopPollingChan:
			log.Debug("[HOST CHECKER] Received kill signal")
			return
		}
	}
}

func (h *HostUptimeChecker) CheckHost(toCheck HostData) {
	log.Debug("[HOST CHECKER] Checking: ", toCheck.CheckURL)

	t1 := time.Now()
	report := HostHealthReport{
		HostData: toCheck,
	}
	switch toCheck.Protocol {
	case "tcp", "tls":
		host := toCheck.CheckURL
		base := toCheck.Protocol + "://"
		if !strings.HasPrefix(host, base) {
			host = base + host
		}
		u, err := url.Parse(host)
		if err != nil {
			log.Error("Could not parse host: ", err)
			return
		}
		var ls net.Conn
		var d net.Dialer
		d.Timeout = toCheck.Timeout
		if toCheck.Protocol == "tls" {
			ls, err = tls.DialWithDialer(&d, "tls", u.Host, nil)
		} else {
			ls, err = d.Dial("tcp", u.Host)
		}
		if err != nil {
			log.Error("Could not connect to host: ", err)
			report.IsTCPError = true
			break
		}
		if toCheck.EnableProxyProtocol {
			log.Debug("using proxy protocol")
			ls = proxyproto.NewConn(ls, 0)
		}
		defer ls.Close()
		for _, cmd := range toCheck.Commands {
			switch cmd.Name {
			case "send":
				log.Debugf("%s: sending %s", host, cmd.Message)
				_, err = ls.Write([]byte(cmd.Message))
				if err != nil {
					log.Errorf("Failed to send %s :%v", cmd.Message, err)
					report.IsTCPError = true
					break
				}
			case "expect":
				buf := make([]byte, len(cmd.Message))
				_, err = ls.Read(buf)
				if err != nil {
					log.Errorf("Failed to read %s :%v", cmd.Message, err)
					report.IsTCPError = true
					break
				}
				g := string(buf)
				if g != cmd.Message {
					log.Errorf("Failed expectation  expected %s got %s", cmd.Message, g)
					report.IsTCPError = true
					break
				}
				log.Debugf("%s: received %s", host, cmd.Message)
			}
		}
		report.ResponseCode = http.StatusOK
	default:
		useMethod := toCheck.Method
		if toCheck.Method == "" {
			useMethod = http.MethodGet
		}
		req, err := http.NewRequest(useMethod, toCheck.CheckURL, strings.NewReader(toCheck.Body))
		if err != nil {
			log.Error("Could not create request: ", err)
			return
		}
		for headerName, headerValue := range toCheck.Headers {
			req.Header.Set(headerName, headerValue)
		}
		req.Header.Set("Connection", "close")
		HostCheckerClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Global().ProxySSLInsecureSkipVerify,
			},
		}
		if toCheck.Timeout != 0 {
			HostCheckerClient.Timeout = toCheck.Timeout
		}
		response, err := HostCheckerClient.Do(req)
		if err != nil {
			report.IsTCPError = true
			break
		}
		response.Body.Close()
		report.ResponseCode = response.StatusCode
	}

	millisec := DurationToMillisecond(time.Since(t1))
	report.Latency = millisec
	if report.IsTCPError {
		h.errorChan <- report
		return
	}

	if report.ResponseCode != http.StatusOK {
		h.errorChan <- report
		return
	}

	// host is healthy, report it
	h.okChan <- report
}

func (h *HostUptimeChecker) Init(workers, triggerLimit, timeout int, hostList map[string]HostData, failureCallback, upCallback, pingCallback func(HostHealthReport)) {

	h.samples = new(sync.Map)
	h.stopPollingChan = make(chan bool)
	h.errorChan = make(chan HostHealthReport)
	h.okChan = make(chan HostHealthReport)
	h.HostList = hostList
	h.failureCallback = failureCallback
	h.upCallback = upCallback
	h.pingCallback = pingCallback

	h.workerPoolSize = workers
	if workers == 0 {
		h.workerPoolSize = defaultWorkerPoolSize
	}

	h.sampleTriggerLimit = triggerLimit
	if triggerLimit == 0 {
		h.sampleTriggerLimit = defaultSampletTriggerLimit
	}

	h.checkTimeout = timeout
	if timeout == 0 {
		h.checkTimeout = defaultTimeout
	}

	log.Debug("[HOST CHECKER] Config:TriggerLimit: ", h.sampleTriggerLimit)
	log.Debug("[HOST CHECKER] Config:Timeout: ~", h.checkTimeout)
	log.Debug("[HOST CHECKER] Config:WorkerPool: ", h.workerPoolSize)

	var err error
	h.pool, err = tunny.CreatePool(h.workerPoolSize, func(hostData interface{}) interface{} {
		input, _ := hostData.(HostData)
		h.CheckHost(input)
		return nil
	}).Open()

	log.Debug("[HOST CHECKER] Init complete")

	if err != nil {
		log.Errorf("[HOST CHECKER POOL] Error: %v\n", err)
	}
}

func (h *HostUptimeChecker) Start() {
	// Start the loop that checks for bum hosts
	h.setStopLoop(false)
	log.Debug("[HOST CHECKER] Starting...")
	go h.HostCheckLoop()
	log.Debug("[HOST CHECKER] Check loop started...")
	go h.HostReporter()
	log.Debug("[HOST CHECKER] Host reporter started...")
}

func (h *HostUptimeChecker) Stop() {
	h.setStopLoop(true)
	h.samples = new(sync.Map)
	h.stopPollingChan <- true
	log.Info("[HOST CHECKER] Stopping poller")
	h.pool.Close()
}

func (h *HostUptimeChecker) ResetList(hostList map[string]HostData) {
	h.resetListMu.Lock()
	h.doResetList = true
	h.newList = hostList
	h.resetListMu.Unlock()
	log.Debug("[HOST CHECKER] Checker reset queued!")
}
