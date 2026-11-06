package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const serviceName = "marketing"

type Metrics struct {
	serviceName string
	provider    *sdkmetric.MeterProvider
}

var httpRequestsTotal metric.Int64Counter

//nolint:exhaustruct_v5
func NewMetrics(
	ctx context.Context,
	serviceName string,
	serviceVersion string,
	otlpEndpoint string,
) (
	*Metrics,
	error,
) {
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otlpEndpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp metric exporter: %w", err)
	}

	res, rscErr := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	))
	if rscErr != nil {
		return nil, fmt.Errorf("failed to create resource: %w", rscErr)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)

	otel.SetMeterProvider(mp)

	m := &Metrics{
		serviceName: serviceName,
		provider:    mp,
	}

	initCustomMetricsErr := m.initCustomMetrics()
	if initCustomMetricsErr != nil {
		return nil, fmt.Errorf(
			"failed to initialize custom metrics: %w",
			initCustomMetricsErr,
		)
	}

	return m, nil
}

func (m *Metrics) initCustomMetrics() error {
	meter := otel.Meter(serviceName)

	var err error

	httpRequestsTotal, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http_requests_total metric: %w", err)
	}

	return nil
}

func IncHTTPRequestTotal() {
	httpRequestsTotal.Add(context.Background(), 1)
}

func (m *Metrics) Shutdown(
	ctx context.Context,
) error {
	err := m.provider.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("metrics shutdown failed: %w", err)
	}

	return nil
}
