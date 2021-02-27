package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type traceware struct {
	service     string
	tracer      oteltrace.Tracer
	propagators propagation.TextMapPropagator
}

// Handler implements the http.Handler interface. It does the actual
// tracing of the request.
func (tw traceware) Handler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		ctx := tw.propagators.Extract(c, toTextMapCarrier{c: c})
		spanName := string(c.RequestURI())
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Method())
		}
		opts := []oteltrace.SpanOption{
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		ctx, span := tw.tracer.Start(ctx, spanName, opts...)
		defer span.End()
		c.SetUserValue(`trace-ctx`, ctx)
		h(c)
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
