package telemetry

import (
	"context"
	"strings"
	"time"

	conf "github.com/horonlee/servora/api/gen/go/conf/v1"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

const shutdownTimeout = 5 * time.Second

func InitTracerProvider(c *conf.Trace, serviceName, env string) (func(), error) {
	endpoint := ""
	if c != nil {
		endpoint = strings.TrimSpace(c.Endpoint)
	}
	if endpoint == "" {
		return func() {}, nil
	}
	if strings.TrimSpace(serviceName) == "" {
		serviceName = "unknown.service"
	}
	if strings.TrimSpace(env) == "" {
		env = "unknown"
	}

	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment.name", env),
			attribute.String("exporter", "otlp"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = tp.Shutdown(ctx)
	}

	return cleanup, nil
}
