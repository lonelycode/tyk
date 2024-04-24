package gateway

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/TykTechnologies/tyk/config"
	"github.com/TykTechnologies/tyk/internal/netutil"
	"github.com/TykTechnologies/tyk/internal/otel"
	"github.com/TykTechnologies/tyk/test"
	"github.com/TykTechnologies/tyk/user"
)

func TestGateway_afterConfSetup(t *testing.T) {

	tests := []struct {
		name           string
		initialConfig  config.Config
		expectedConfig config.Config
	}{
		{
			name: "slave options test",
			initialConfig: config.Config{
				SlaveOptions: config.SlaveOptionsConfig{
					UseRPC: true,
				},
			},
			expectedConfig: config.Config{
				SlaveOptions: config.SlaveOptionsConfig{
					UseRPC:                   true,
					GroupID:                  "ungrouped",
					CallTimeout:              30,
					PingTimeout:              60,
					KeySpaceSyncInterval:     10,
					RPCCertCacheExpiration:   3600,
					RPCGlobalCacheExpiration: 30,
				},
				AnalyticsConfig: config.AnalyticsConfigConfig{
					PurgeInterval: 10,
				},
				HealthCheckEndpointName: "hello",
			},
		},
		{
			name: "opentelemetry options test",
			initialConfig: config.Config{
				OpenTelemetry: otel.OpenTelemetry{
					Enabled: true,
				},
			},
			expectedConfig: config.Config{
				OpenTelemetry: otel.OpenTelemetry{
					Enabled:            true,
					Exporter:           "grpc",
					Endpoint:           "localhost:4317",
					ResourceName:       "tyk-gateway",
					SpanProcessorType:  "batch",
					ConnectionTimeout:  1,
					ContextPropagation: "tracecontext",
					Sampling: otel.Sampling{
						Type: "AlwaysOn",
					},
				},
				AnalyticsConfig: config.AnalyticsConfigConfig{
					PurgeInterval: 10,
				},
				HealthCheckEndpointName: "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gw := NewGateway(tt.initialConfig, context.Background())
			gw.afterConfSetup()

			assert.Equal(t, tt.expectedConfig, gw.GetConfig())

		})
	}
}

func TestGateway_apisByIDLen(t *testing.T) {
	tcs := []struct {
		name     string
		APIs     []string
		expected int
	}{
		{
			name:     "empty apis",
			APIs:     []string{},
			expected: 0,
		},
		{
			name:     "one api",
			APIs:     []string{"api1"},
			expected: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ts := StartTest(nil)
			defer ts.Close()

			for i := range tc.APIs {
				ts.Gw.BuildAndLoadAPI(func(spec *APISpec) {
					spec.APIID = tc.APIs[i]
					spec.UseKeylessAccess = false
					spec.OrgID = "default"
				})
			}

			actual := ts.Gw.apisByIDLen()

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGateway_policiesByIDLen(t *testing.T) {
	tcs := []struct {
		name     string
		policies []string
		expected int
	}{
		{
			name:     "empty policies",
			policies: []string{},
			expected: 0,
		},
		{
			name:     "one policy",
			policies: []string{"policy1"},
			expected: 1,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ts := StartTest(nil)
			defer ts.Close()

			for _, pol := range tc.policies {
				ts.CreatePolicy(func(p *user.Policy) {
					p.Name = pol
				})
			}

			actual := ts.Gw.policiesByIDLen()

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGateway_SyncResourcesWithReload(t *testing.T) {
	retryAttempts := 2
	ts := StartTest(func(globalConf *config.Config) {
		globalConf.ResourceSync.RetryAttempts = retryAttempts
		globalConf.ResourceSync.Interval = 1
	})

	var syncErr = errors.New("sync error")
	syncFuncSuccessAt := func(t *testing.T, successAt int) (func() (int, error), *int) {
		t.Helper()
		var hitCount int
		return func() (int, error) {
			hitCount++
			if hitCount == successAt {
				return 10, nil
			}
			return 0, syncErr
		}, &hitCount
	}

	t.Run("invalid resource", func(t *testing.T) {
		t.Parallel()
		syncFunc, hitCounter := syncFuncSuccessAt(t, 0)
		resourceCount, err := syncResourcesWithReload("unknown-resource", ts.Gw.GetConfig(), syncFunc)
		assert.Error(t, ErrSyncResourceNotKnown, err)
		assert.Zero(t, resourceCount)
		assert.Zero(t, *hitCounter)
	})

	t.Run("sync success at first try", func(t *testing.T) {
		t.Parallel()
		syncFunc, hitCounter := syncFuncSuccessAt(t, 1)
		resourceCount, err := syncResourcesWithReload("apis", ts.Gw.GetConfig(), syncFunc)
		assert.NoError(t, err)
		assert.Equal(t, 10, resourceCount)
		assert.Equal(t, 1, *hitCounter)
	})

	t.Run("sync failed after retries", func(t *testing.T) {
		t.Parallel()
		syncFunc, hitCounter := syncFuncSuccessAt(t, 5)
		startTime := time.Now()
		resourceCount, err := syncResourcesWithReload("apis", ts.Gw.GetConfig(), syncFunc)
		assert.Greater(t, time.Since(startTime), time.Second*3)
		assert.ErrorIs(t, err, syncErr)
		assert.Zero(t, resourceCount)
		assert.Equal(t, 3, *hitCounter)
	})

	t.Run("sync success after first retry", func(t *testing.T) {
		t.Parallel()
		syncFunc, hitCounter := syncFuncSuccessAt(t, 2)
		startTime := time.Now()
		resourceCount, err := syncResourcesWithReload("apis", ts.Gw.GetConfig(), syncFunc)
		assert.Greater(t, time.Since(startTime), time.Second*1)
		assert.NoError(t, err)
		assert.Equal(t, 10, resourceCount)
		assert.Equal(t, 2, *hitCounter)
	})

}

type gatewayGetHostDetailsTestCheckFn func(*testing.T, *test.BufferedLogger, *Gateway)

func gatewayGetHostDetailsTestHasErr(wantErr bool, errorText string) gatewayGetHostDetailsTestCheckFn {
	return func(t *testing.T, bl *test.BufferedLogger, _ *Gateway) {
		logs := bl.GetLogs(logrus.ErrorLevel)
		if !wantErr && assert.Empty(t, logs) {
			return
		}

		if wantErr && !assert.NotEmpty(t, logs) {
			return
		}

		if wantErr && errorText != "" {
			for _, log := range logs {
				assert.Contains(t, log.Message, errorText)
			}
		}
	}
}

func gatewayGetHostDetailsTesHasAddress(addr string) gatewayGetHostDetailsTestCheckFn {
	return func(t *testing.T, bl *test.BufferedLogger, gw *Gateway) {
		matched, err := regexp.MatchString(addr, gw.hostDetails.Address)
		if err != nil {
			t.Errorf("Failed to compile regex pattern: %v", err)
		}
		if !matched {
			t.Errorf("Wanted address %s, got %s", addr, gw.hostDetails.Address)
		}
	}
}

func defineGatewayGetHostDetailsTests() []struct {
	name                string
	before              func(*Gateway)
	readPIDFromFile     func(string) (int, error)
	netutilGetIpAddress func() ([]string, error)
	checks              []gatewayGetHostDetailsTestCheckFn
	// config              config.Config
} {
	var (
		check = func(fns ...gatewayGetHostDetailsTestCheckFn) []gatewayGetHostDetailsTestCheckFn { return fns }
		// matches ipv6 and ipv4
		// https://go.dev/play/p/cbOQUiqNqyU
		// https://regex101.com/r/lIWfaA/1
		ipAddrPattern = `^((([a-f0-9]{1,4}):){7}([a-f0-9]{1,4})|(([a-f0-9]{1,4})(:([a-f0-9]{1,4})){0,6})?::(([a-f0-9]{1,4})(:([a-f0-9]{1,4})){0,6})?|((([a-f0-9]{1,4}):){5}[a-f0-9]{1,4}|(([a-f0-9]{1,4}):){0,5}:([a-f0-9]{1,4}))|(([a-f0-9]{1,4}:){0,6}[a-f0-9]{1,4}|(([a-f0-9]{1,4}:){0,6}:([a-f0-9]{1,4})){0,1})|([a-f0-9]{1,4}:){0,7}:|([a-f0-9]{1,4}:){0,6}[a-f0-9]{1,4}|(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))$`
	)
	return []struct {
		name                string
		before              func(*Gateway)
		readPIDFromFile     func(string) (int, error)
		netutilGetIpAddress func() ([]string, error)
		checks              []gatewayGetHostDetailsTestCheckFn
		// config              config.Config
	}{
		{
			name:            "fail-read-pid",
			readPIDFromFile: func(file string) (int, error) { return 0, fmt.Errorf("Error opening file") },
			before: func(gw *Gateway) {
				gw.SetConfig(config.Config{
					ListenAddress: "127.0.0.1",
				})
			},
			checks: check(
				gatewayGetHostDetailsTestHasErr(true, "Error opening file"),
			),
		},
		{
			name:            "success-listen-address-set",
			readPIDFromFile: func(file string) (int, error) { return 1000, nil },
			before: func(gw *Gateway) {
				gw.SetConfig(config.Config{
					ListenAddress: "127.0.0.1",
				})
			},
			checks: check(
				gatewayGetHostDetailsTestHasErr(false, ""),
				gatewayGetHostDetailsTesHasAddress("127\\.0\\.0\\.1"),
			),
		},
		{
			name:            "success-listen-address-not-set",
			readPIDFromFile: func(file string) (int, error) { return 1000, nil },
			before: func(gw *Gateway) {
				gw.SetConfig(config.Config{
					ListenAddress: "",
				})
			},
			checks: check(
				gatewayGetHostDetailsTestHasErr(false, ""),
				gatewayGetHostDetailsTesHasAddress(ipAddrPattern),
			),
		},
		{
			name:            "fail-getting-network-address",
			readPIDFromFile: func(file string) (int, error) { return 1000, nil },
			before: func(gw *Gateway) {
				gw.SetConfig(config.Config{
					ListenAddress: "",
				})
			},
			netutilGetIpAddress: func() ([]string, error) { return nil, fmt.Errorf("Error getting network addresses") },
			checks: check(
				gatewayGetHostDetailsTestHasErr(true, "Error getting network addresses"),
			),
		}, // Define your test cases here
	}
}

func TestGatewayGetHostDetails(t *testing.T) {

	var (
		orig_readPIDFromFile = readPIDFromFile
		orig_mainLog         = mainLog
		orig_getIpAddress    = netutil.GetIpAddress
		bl                   = test.NewBufferingLogger()
	)

	tests := defineGatewayGetHostDetailsTests()

	// restore the original functions
	defer func() {
		readPIDFromFile = orig_readPIDFromFile
		mainLog = orig_mainLog
		getIpAddress = orig_getIpAddress
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// clear logger mock buffer
			bl.ClearLogs()
			// replace fucntions with mocks
			mainLog = bl.Logger.WithField("prefix", "test")
			if tt.readPIDFromFile != nil {
				readPIDFromFile = tt.readPIDFromFile
			}

			if tt.netutilGetIpAddress != nil {
				getIpAddress = tt.netutilGetIpAddress
			}

			gw := &Gateway{}

			if tt.before != nil {
				tt.before(gw)
			}

			gw.getHostDetails(gw.GetConfig().PIDFileLocation)
			for _, c := range tt.checks {
				c(t, bl, gw)
			}
		})
	}
}
