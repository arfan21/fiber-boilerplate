package telemetry

import (
	"context"

	"github.com/arfan21/fiber-boilerplate/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitLogs() (func(context.Context) error, error) {
	ctx := context.Background()
	secureOption := otlploggrpc.WithInsecure()
	grpcExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(collectorURL),
		secureOption,
	)

	if err != nil {
		logger.Log(ctx).Error().Err(err).Msg("logs: could not initialize grpc client")
		return nil, err
	}

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Log(ctx).Error().Err(err).Msg("logs: could not initialize resource")
		return nil, err
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(grpcExporter)),
		sdklog.WithResource(resources),
	)

	global.SetLoggerProvider(provider)

	return provider.Shutdown, nil
}
