package initializers

import (
	"context"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
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

	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
)

func NewMetricsAndTracer(v *viper.Viper) (func() metric.MeterProvider, func() *sdktrace.TracerProvider) {
	ctx := context.Background()
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(v.GetString(`OTLP_GRPC`)),
	)
	exp, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		panic(err)

	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(v.GetString(`SERVICE_NAME`)),
		),
	)
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.TraceIDRatioBased(v.GetFloat64(`TRACE_RATIO`))}),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	if err != nil {
		panic(err)
	}

	cont := controller.New(
		processor.New(
			simple.NewWithHistogramDistribution(),
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

	return cont.MeterProvider, func() *sdktrace.TracerProvider {
		return tracerProvider
	}
}
