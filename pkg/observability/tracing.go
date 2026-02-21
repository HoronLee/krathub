package observability

import (
	"context"
	"time"

	conf "github.com/horonlee/krathub/api/gen/go/conf/v1"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const shutdownTimeout = 5 * time.Second

func InitTracerProvider(c *conf.Trace, serviceName, env string) (func(), error) {
	if c == nil || c.Endpoint == "" {
		return func() {}, nil
	}

	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(c.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("exporter", "otlp"),
			attribute.String("env", env),
		)),
	)
	otel.SetTracerProvider(tp)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = tp.Shutdown(ctx)
	}

	return cleanup, nil
}
