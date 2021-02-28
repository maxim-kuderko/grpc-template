package main

import (
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv"
	oteltrace "go.opentelemetry.io/otel/trace"
	"time"
)

type traceware struct {
	service     string
	tracer      oteltrace.Tracer
	propagators propagation.TextMapPropagator
}

func (tw traceware) Handler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	opts := []oteltrace.SpanOption{
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
	}
	m, err := global.GetMeterProvider().Meter(`fastthttp`).NewInt64ValueRecorder(`response-time`)
	if err != nil {
		panic(err)
	}
	return func(c *fasthttp.RequestCtx) {
		ctx := tw.propagators.Extract(c, toTextMapCarrier{c: c})
		spanName := string(c.Request.URI().Path())
		ctx, span := tw.tracer.Start(ctx, spanName, opts...)
		defer span.End()
		c.SetUserValue(`trace-ctx`, ctx)
		t := time.Now()
		h(c)
		m.Record(c, time.Since(t).Microseconds(), label.Any(`handler`, string(c.Request.URI().Path())))
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(c.Response.StatusCode())
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(c.Response.StatusCode())
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
	}
}

type toTextMapCarrier struct {
	c *fasthttp.RequestCtx
}

func (t toTextMapCarrier) Get(key string) string {
	return string(t.c.Request.Header.Peek(key))
}

func (t toTextMapCarrier) Set(key string, value string) {
	t.c.Response.Header.Set(key, value)
}
