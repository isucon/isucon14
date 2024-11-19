package metrics

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/isucon/isucon14/bench/benchrun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/noop"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func Setup(noOp bool) (metricsdk.Exporter, error) {
	if noOp {
		otel.SetMeterProvider(noop.NewMeterProvider())
		return stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
	}

	exp, err := getExporter()
	if err != nil {
		return nil, err
	}

	recources, err := resource.New(
		context.Background(),
		resource.WithProcessCommandArgs(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(benchrun.GetTargetAddress()),
		),
	)
	if err != nil {
		return nil, err
	}

	reader := metricsdk.NewPeriodicReader(exp, metricsdk.WithInterval(3*time.Second))
	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(recources),
		metricsdk.WithReader(reader),
	)
	otel.SetMeterProvider(provider)

	return exp, nil
}

func getExporter() (metricsdk.Exporter, error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		return otlpmetrichttp.New(context.Background())
	}
	return stdoutmetric.New(stdoutmetric.WithWriter(os.Stderr))
}
