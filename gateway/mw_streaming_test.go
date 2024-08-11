package gateway

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"

	"github.com/IBM/sarama"
	"github.com/testcontainers/testcontainers-go/modules/kafka"

	"github.com/TykTechnologies/tyk/apidef/oas"
	"github.com/TykTechnologies/tyk/config"
	"github.com/TykTechnologies/tyk/test"
)

// ConvertYAMLToJSON converts a YAML byte slice to a JSON byte slice
func ConvertYAMLToJSON(yamlData []byte) ([]byte, error) {
	var rawData interface{}
	if err := yaml.Unmarshal(yamlData, &rawData); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	// Convert map[interface{}]interface{} to map[string]interface{}
	data := convertToStringKeyMap(rawData)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling to JSON: %w", err)
	}

	return jsonData, nil
}

// convertToStringKeyMap recursively converts map[interface{}]interface{} to map[string]interface{}
func convertToStringKeyMap(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range x {
			m[fmt.Sprintf("%v", k)] = convertToStringKeyMap(v)
		}
		return m
	case []interface{}:
		for i, v := range x {
			x[i] = convertToStringKeyMap(v)
		}
	}
	return i
}

func TestAsyncAPI(t *testing.T) {
	t.SkipNow()

	ts := StartTest(func(globalConf *config.Config) {
		globalConf.Labs = map[string]interface{}{
			"streaming": map[string]interface{}{
				"enabled": true,
			},
		}
	})

	ts.Gw.BuildAndLoadAPI(func(spec *APISpec) {
		spec.Proxy.ListenPath = "/test"
		spec.UseKeylessAccess = true
	})

	// Check that standard API works
	_, _ = ts.Run(t, test.TestCase{Code: http.StatusOK, Method: http.MethodGet, Path: "/test"})

	defer ts.Close()

	const (
		oldAPIID    = "old-api-id"
		oasAPIID    = "oas-api-id"
		oasBasePath = "/tyk/apis/oas"
	)

	tempFile, err := os.CreateTemp("", "test-output-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name()) // clean up
	tempFilePath := tempFile.Name()

	streamingConfig := `
streams:
  test:
    input:
      label: "in"
      generate:
        count: 3
        interval: ""
        mapping: root = "hello world"

    output:
      label: "out"
      file:
        path: "%s"
        codec: lines`

	streamingConfig = fmt.Sprintf(streamingConfig, tempFilePath)

	streamingConfigJSON, err := ConvertYAMLToJSON([]byte(streamingConfig))
	if err != nil {
		t.Fatalf("Failed to convert YAML to JSON: %v", err)
	}

	var parsedStreamingConfig map[string]interface{}
	if err := json.Unmarshal(streamingConfigJSON, &parsedStreamingConfig); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	oasAPI := oas.OAS{
		T: openapi3.T{
			OpenAPI: "3.0.3",
			Info: &openapi3.Info{
				Title:   "oas doc",
				Version: "1",
			},
			Paths: make(openapi3.Paths),
		},
	}

	oasAPI.Extensions = map[string]interface{}{
		ExtensionTykStreaming: parsedStreamingConfig,
		// oas.ExtensionTykAPIGateway: tykExtension,
	}

	ts.Gw.BuildAndLoadAPI(func(spec *APISpec) {
		spec.Proxy.ListenPath = "/test"
		spec.UseKeylessAccess = true
	}, func(spec *APISpec) {
		spec.SetDisabledFlags()
		spec.APIID = "base-api-id"
		spec.VersionDefinition.Enabled = false
		spec.VersionDefinition.Key = ""
		spec.VersionDefinition.Location = ""

		spec.IsOAS = true
		spec.OAS = oasAPI
		spec.OAS.Fill(*spec.APIDefinition)
	})

	// Check that standard API still works
	_, _ = ts.Run(t, test.TestCase{Code: http.StatusOK, Method: http.MethodGet, Path: "/test"})

	if globalStreamCounter.Load() != 1 {
		t.Fatalf("Expected 1 stream, got %d", globalStreamCounter.Load())
	}

	time.Sleep(500 * time.Millisecond)

	content, err := os.ReadFile(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	lines := strings.Split(string(content), "\n")

	// Adjust for the trailing new line which results in an extra empty element in the slice
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}

	for _, line := range lines {
		if line != "hello world" {
			t.Fatalf("Expected 'hello world', got '%s'", line)
		}
	}
}

var tests = []struct {
	name          string
	consumerGroup string
	tenantID      string
	isDynamic     bool
}{
	{"StaticGroup", "static-group", "default", false},
	{"DynamicGroup", "$tyk_context.request_id", "dynamic", true},
}

func TestAsyncAPIHttp(t *testing.T) {
	ctx := context.Background()
	kafkaContainer, err := kafka.Run(ctx, "confluentinc/confluent-local:7.5.0")
	if err != nil {
		t.Fatalf("Failed to start Kafka container: %v", err)
	}
	defer kafkaContainer.Terminate(ctx)

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		t.Fatalf("Failed to get Kafka brokers: %v", err)
	}

	ts := StartTest(func(globalConf *config.Config) {
		globalConf.Labs = map[string]interface{}{
			"streaming": map[string]interface{}{
				"enabled": true,
			},
		}
	})
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			apiName := setupStreamingAPI(t, ts, tc.consumerGroup, tc.tenantID, brokers[0])
			testAsyncAPIHttp(t, ts, tc.consumerGroup, tc.isDynamic, tc.tenantID, apiName, brokers[0])
		})
	}
}

func setupStreamingAPI(t *testing.T, ts *Test, consumerGroup string, tenantID string, kafkaHost string) string {
	t.Logf("Setting up streaming API for tenant: %s with consumer group: %s", tenantID, consumerGroup)

	apiName := fmt.Sprintf("streaming-api-%s", tenantID)

	ts.Gw.BuildAndLoadAPI(func(spec *APISpec) {
		spec.Proxy.ListenPath = fmt.Sprintf("/%s", apiName)
		spec.UseKeylessAccess = true
		spec.IsOAS = true
		spec.OAS = setupOASForStreamingAPI(t, consumerGroup, kafkaHost)
		spec.OAS.Fill(*spec.APIDefinition)
	})

	return apiName
}

func setupOASForStreamingAPI(t *testing.T, consumerGroup string, kafkaHost string) oas.OAS {
	streamingConfig := fmt.Sprintf(`
streams:
 test:
  input:
   kafka:
    addresses: ["%s"]
    topics: ["test"]
    consumer_group: "%s"

  output:
   http_server:
    path: /get
    ws_path: /get/ws`, kafkaHost, consumerGroup)

	streamingConfigJSON, err := ConvertYAMLToJSON([]byte(streamingConfig))
	if err != nil {
		t.Fatalf("Failed to convert YAML to JSON: %v", err)
	}

	var parsedStreamingConfig map[string]interface{}
	if err := json.Unmarshal(streamingConfigJSON, &parsedStreamingConfig); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	oasAPI := oas.OAS{
		T: openapi3.T{
			OpenAPI: "3.0.3",
			Info: &openapi3.Info{
				Title:   "oas doc",
				Version: "1",
			},
			Paths: make(openapi3.Paths),
		},
	}

	oasAPI.Extensions = map[string]interface{}{
		ExtensionTykStreaming: parsedStreamingConfig,
	}

	return oasAPI
}

func testAsyncAPIHttp(t *testing.T, ts *Test, consumerGroup string, isDynamic bool, tenantID string, apiName string, kafkaHost string) {
	const messageToSend = "hello websocket"
	const numMessages = 2
	const numClients = 2

	streamCount := globalStreamCounter.Load()
	t.Logf("Stream count for tenant %s: %d", tenantID, streamCount)

	// Create WebSocket clients
	wsClients := make([]*websocket.Conn, numClients)
	for i := 0; i < numClients; i++ {
		dialer := websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		}
		wsURL := strings.Replace(ts.URL, "http", "ws", 1) + fmt.Sprintf("/%s/get/ws", apiName)
		wsConn, resp, err := dialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Failed to connect to WebSocket %d: %v\nResponse: %+v", i+1, err, resp)
		}
		defer wsConn.Close()
		wsClients[i] = wsConn
		t.Logf("Successfully connected to WebSocket %d. Response: %+v", i+1, resp)
	}

	// Add a delay to ensure WebSocket connections are fully established
	time.Sleep(1 * time.Second)

	// Send messages to Kafka
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{kafkaHost}, config)
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	for i := 0; i < numMessages; i++ {
		msg := &sarama.ProducerMessage{
			Topic: "test",
			Value: sarama.StringEncoder(fmt.Sprintf("%s %d", messageToSend, i+1)),
		}

		t.Logf("Sending message to Kafka topic 'test': %s", msg.Value)
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			t.Fatalf("Failed to send message to Kafka: %v", err)
		}
		t.Logf("Message sent to partition %d at offset %d", partition, offset)
	}

	expectedTotalMessages := numMessages
	if isDynamic {
		expectedTotalMessages *= numClients
	}

	messagesReceived := 0
	overallTimeout := time.After(10 * time.Second)
	inactivityTimeout := time.NewTimer(2 * time.Second)
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic in WebSocket read goroutine: %v", r)
				done <- true
			}
		}()

		for {
			select {
			case <-overallTimeout:
				t.Log("Overall timeout reached while waiting for messages")
				done <- true
				return
			case <-inactivityTimeout.C:
				t.Log("Inactivity timeout reached")
				done <- true
				return
			default:
				for i, wsConn := range wsClients {
					wsConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
					_, p, err := wsConn.ReadMessage()
					if err != nil {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
							t.Logf("Unexpected error reading from WebSocket %d: %v", i+1, err)
						} else if !strings.Contains(err.Error(), "i/o timeout") {
							t.Logf("Error reading from WebSocket %d: %v", i+1, err)
						}
					} else {
						receivedMessage := string(p)
						t.Logf("Received message from WebSocket %d: %s", i+1, receivedMessage)
						if strings.HasPrefix(receivedMessage, messageToSend) {
							messagesReceived++
							t.Logf("Message from WebSocket %d matches sent message", i+1)
							inactivityTimeout.Reset(2 * time.Second)
						}
					}
				}

				if messagesReceived >= expectedTotalMessages {
					t.Logf("Received all expected messages (%d)", messagesReceived)
					done <- true
					return
				}

				// Add a small sleep to prevent tight loop
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	<-done

	if messagesReceived != expectedTotalMessages {
		t.Errorf("Expected %d messages, but received %d for tenant %s", expectedTotalMessages, messagesReceived, tenantID)
	} else {
		t.Logf("Successfully received %d messages as expected for tenant %s", messagesReceived, tenantID)
	}
}

func waitForAPIToBeLoaded(ts *Test) error {
	maxAttempts := 2
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(ts.URL + "/streaming-api-default/metrics")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("API failed to load after %d attempts", maxAttempts)
}

func TestWebSocketConnectionClosedOnAPIReload(t *testing.T) {
	ctx := context.Background()
	kafkaContainer, err := kafka.Run(ctx, "confluentinc/confluent-local:7.5.0")
	if err != nil {
		t.Fatalf("Failed to start Kafka container: %v", err)
	}
	defer kafkaContainer.Terminate(ctx)

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		t.Fatalf("Failed to get Kafka brokers: %v", err)
	}

	ts := StartTest(func(globalConf *config.Config) {
		globalConf.Labs = map[string]interface{}{
			"streaming": map[string]interface{}{
				"enabled": true,
			},
		}
	})
	defer ts.Close()

	apiName := setupStreamingAPI(t, ts, "test-group", "default", brokers[0])

	// Connect to WebSocket
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}
	wsURL := strings.Replace(ts.URL, "http", "ws", 1) + fmt.Sprintf("/%s/get/ws", apiName)
	wsConn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer wsConn.Close()

	// Reload the API by rebuilding and loading it
	ts.Gw.BuildAndLoadAPI(func(spec *APISpec) {
		spec.Proxy.ListenPath = fmt.Sprintf("/%s", apiName)
		spec.UseKeylessAccess = true
		spec.IsOAS = true
		spec.OAS = setupOASForStreamingAPI(t, "test-group", brokers[0])
	})

	// Wait for the API to be reloaded
	err = waitForAPIToBeLoaded(ts)
	if err != nil {
		t.Fatalf("API failed to reload: %v", err)
	}

	// Try to send a message, which should fail if the connection is closed
	err = wsConn.WriteMessage(websocket.TextMessage, []byte("test message"))
	if err == nil {
		t.Fatalf("Expected WebSocket connection to be closed, but write succeeded")
	}

	// Verify that the error indicates a closed connection
	if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		t.Fatalf("Expected WebSocket to be closed with CloseGoingAway or CloseAbnormalClosure, got: %v", err)
	}

	t.Log("WebSocket connection was successfully closed on API reload")
}
