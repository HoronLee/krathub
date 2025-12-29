package server

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Metrics is a metrics struct.
type Metrics struct {
	Requests metric.Int64Counter
	Seconds  metric.Float64Histogram
	Handler  http.Handler
}

// NewMetrics initializes metrics.
func NewMetrics(c *conf.Metrics, logger log.Logger) (*Metrics, error) {
	if c == nil || !c.Enable {
		log.NewHelper(logger).Info("metrics config is empty, skip metrics init")
		return nil, nil
	}

	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))

	meterName := "krathub"
	if c.MeterName != "" {
		meterName = c.MeterName
	}
	meter := provider.Meter(meterName)

	requests, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		return nil, err
	}

	seconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Requests: requests,
		Seconds:  seconds,
		Handler:  promhttp.Handler(),
	}, nil
}
