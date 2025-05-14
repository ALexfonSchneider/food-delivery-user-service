package otel

import (
	"context"
	"fmt"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/config"
	"github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"google.golang.org/grpc"
	"log/slog"
	"os"
)

func InitTracerProvider(conf *config.Config) (*sdktrace.TracerProvider, func()) {
	ctx := context.Background()

	// Подключение к OTLP endpoint (например, Otel Collector или напрямую в Tempo/Jaeger)
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(conf.OpenTelemetry.GRPCEndpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		panic(err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String("user-service")),
	)
	if err != nil {
		panic(err)
	}

	// Создание провайдера трейсов
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)

	// Регистрация глобального провайдера
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint(conf.OpenTelemetry.GRPCEndpoint), otlploggrpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize OTLP log exporter: %v", err))
	}

	lp := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)

	global.SetLoggerProvider(lp)

	logger := slog.New(
		slogmulti.Fanout(
			otelslog.NewHandler("user-service"),
			slog.NewJSONHandler(os.Stdout, nil),
		),
	)

	slog.SetDefault(logger)

	return tp, func() {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("error shutting down tracer provider: %v", err)
		}
		if err := lp.Shutdown(ctx); err != nil {
			panic(fmt.Sprintf("failed to shutdown logger provider: %v", err))
		}
	}
}
