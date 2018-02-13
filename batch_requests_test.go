package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TykTechnologies/tyk/apidef"
	"github.com/TykTechnologies/tyk/config"
	"github.com/TykTechnologies/tyk/test"
)

const testBatchRequest = `{
	"requests": [
	{
		"method": "GET",
		"headers": {
			"test-header-1": "test-1",
			"test-header-2": "test-2"
		},
		"relative_url": "get/?param1=this"
	},
	{
		"method": "POST",
		"body": "TEST BODY",
		"relative_url": "post/"
	},
	{
		"method": "PUT",
		"relative_url": "put/"
	}
	],
	"suppress_parallel_execution": true
}`

func TestBatch(t *testing.T) {
	ts := newTykTestServer()
	defer ts.Close()

	buildAndLoadAPI(func(spec *APISpec) {
		spec.Proxy.ListenPath = "/v1/"
		spec.EnableBatchRequestSupport = true
	})

	ts.Run(t, []test.TestCase{
		{Method: "POST", Path: "/v1/tyk/batch/", Data: `{"requests":[]}`, Code: 200, BodyMatch: "[]"},
		{Method: "POST", Path: "/v1/tyk/batch/", Data: "malformed", Code: 400},
		{Method: "POST", Path: "/v1/tyk/batch/", Data: testBatchRequest, Code: 200},
	}...)

	resp, _ := ts.Do(test.TestCase{Method: "POST", Path: "/v1/tyk/batch/", Data: testBatchRequest})
	if resp != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		var batchResponse []map[string]json.RawMessage
		if err := json.Unmarshal(body, &batchResponse); err != nil {
			t.Fatal(err)
		}

		if len(batchResponse) != 3 {
			t.Errorf("Length not match: %d", len(batchResponse))
		}

		if string(batchResponse[0]["relative_url"]) != `"get/?param1=this"` {
			t.Error("Url order not match:", string(batchResponse[0]["relative_url"]))
		}
	}
}

var virtBatchTest = `function batchTest (request, session, config) {
    // Set up a response object
    var response = {
        Body: "",
        Headers: {
            "content-type": "application/json"
        },
        Code: 202
    }
    
    // Batch request
    var batch = {
        "requests": [
			{
				"method": "GET",
				"headers": {},
                "body": "",
				"relative_url": "{upstream_URL}"
			},
			{
				"method": "GET",
				"headers": {},
                "body": "",
				"relative_url": "{upstream_URL}"
			}
		],
        "suppress_parallel_execution": false
    }
	
    var newBody = TykBatchRequest(JSON.stringify(batch))
	var asJS = JSON.parse(newBody)
	for (var i in asJS) {
		if (asJS[i].code == 0){
			response.Code = 404
		}	
    }
    return TykJsResponse(response, session.meta_data)
	
}`

func TestSSLBatch(t *testing.T) {

	_, _, combinedClientPEM, clientCert := genCertificate(&x509.Certificate{})
	clientCert.Leaf, _ = x509.ParseCertificate(clientCert.Certificate[0])

	upstream := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	// Mutual TLS protected upstream
	pool := x509.NewCertPool()
	pool.AddCert(clientCert.Leaf)
	upstream.TLS = &tls.Config{
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          pool,
		InsecureSkipVerify: true,
	}

	upstream.StartTLS()
	defer upstream.Close()

	_, _, combinedPEM, _ := genServerCertificate()
	serverCertID, _ := CertificateManager.Add(combinedPEM, "")
	defer CertificateManager.Delete(serverCertID)

	clientCertID, _ := CertificateManager.Add(combinedClientPEM, "")
	defer CertificateManager.Delete(clientCertID)

	virtBatchTest = strings.Replace(virtBatchTest, "{upstream_URL}", upstream.URL, 2)
	defer upstream.Close()
	log.Debug(upstream.URL)
	upstreamHost := strings.TrimPrefix(upstream.URL, "https://")

	config.Global.Security.Certificates.Upstream = map[string]string{upstreamHost: clientCertID}
	config.Global.HttpServerOptions.UseSSL = true
	config.Global.HttpServerOptions.SSLCertificates = []string{serverCertID}
	config.Global.ProxySSLInsecureSkipVerify = true

	defer resetTestConfig()

	ts := newTykTestServer()
	defer ts.Close()

	buildAndLoadAPI(func(spec *APISpec) {
		spec.Proxy.ListenPath = "/"
		virtualMeta := apidef.VirtualMeta{
			ResponseFunctionName: "batchTest",
			FunctionSourceType:   "blob",
			FunctionSourceURI:    base64.StdEncoding.EncodeToString([]byte(virtBatchTest)),
			Path:                 "/virt",
			Method:               "GET",
		}
		v := spec.VersionData.Versions["v1"]
		v.UseExtendedPaths = true
		v.ExtendedPaths = apidef.ExtendedPathsSet{
			Virtual: []apidef.VirtualMeta{virtualMeta},
		}
		spec.VersionData.Versions["v1"] = v
	})
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	}}}

	t.Run("Skip verification", func(t *testing.T) {
		ts.Run(t, test.TestCase{
			Path: "/virt", Code: 202, Client: client,
		})
	})

	t.Run("Verification required", func(t *testing.T) {

		config.Global.ProxySSLInsecureSkipVerify = false

		ts.Run(t, test.TestCase{
			Path: "/virt", Code: 404, Client: client,
		})
	})

}
