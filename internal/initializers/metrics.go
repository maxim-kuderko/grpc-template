package initializers

import (
	"context"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"time"
)

func NewMetrics(v *viper.Viper) func() metric.MeterProvider {
	ctx := context.Background()
	exp, _ := otlp.NewExporter(ctx, otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(v.GetString(`OTLP_GRPC`)),
	))
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(v.GetString(`SERVICE_NAME`)),
		),
	)
	if err != nil {
		panic(err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.TraceIDRatioBased(v.GetFloat64(`TRACE_RATIO`))}),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	cont := controller.New(
		processor.New(
			simple.NewWithInexpensiveDistribution(),
			exp,
		),
		controller.WithPusher(exp),
		controller.WithCollectPeriod(2*time.Second),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)
	global.SetMeterProvider(cont.MeterProvider())
	if err := cont.Start(context.Background()); err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return cont.MeterProvider
}
