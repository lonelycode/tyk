package gateway

import (
	"github.com/sirupsen/logrus"

	tyktrace "github.com/TykTechnologies/opentelemetry/trace"
)

// initOpenTelemetry initializes OpenTelemetry tracer provider
// in case of invalid configuration or disabled OpenTelemetry, it will initialize a noop tracer
func (gw *Gateway) initOpenTelemetry() {
	gwConfig := gw.GetConfig()

	traceLogger := mainLog.WithFields(logrus.Fields{
		"exporter":           gwConfig.OpenTelemetry.Exporter,
		"endpoint":           gwConfig.OpenTelemetry.Endpoint,
		"connection_timeout": gwConfig.OpenTelemetry.ConnectionTimeout,
	})

	var errOtel error
	gw.TraceProvider, errOtel = tyktrace.NewProvider(
		tyktrace.WithContext(gw.ctx),
		tyktrace.WithConfig(&gwConfig.OpenTelemetry),
		tyktrace.WithLogger(traceLogger),
	)

	if errOtel != nil {
		mainLog.Errorf("Initializing OpenTelemetry %s", errOtel)
	}
}
