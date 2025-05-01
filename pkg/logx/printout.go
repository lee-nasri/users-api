package logx

import (
	"context"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Infow(ctx context.Context, message string, kvs ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)

	infow(message, traceInfo, kvs...)
}

func infow(message string, trace DatadogTracing, kvs ...interface{}) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Sugar().Infow(message, kvs...)
}

func Errorw(ctx context.Context, message string, kvs ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)

	errorw(message, traceInfo, kvs...)
}

func errorw(message string, trace DatadogTracing, kvs ...interface{}) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Sugar().Errorw(message, kvs...)
}

func Warnw(ctx context.Context, message string, kvs ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)

	warnw(message, traceInfo, kvs...)
}

func warnw(message string, trace DatadogTracing, kvs ...interface{}) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Sugar().Warnw(message, kvs...)
}

func Debugw(ctx context.Context, message string, kvs ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)

	debugw(message, traceInfo, kvs...)
}

func debugw(message string, trace DatadogTracing, kvs ...interface{}) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Sugar().Debugw(message, kvs...)
}

func Fatalw(ctx context.Context, message string, kvs ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)

	fatalw(message, traceInfo, kvs...)
}

func fatalw(message string, trace DatadogTracing, kvs ...interface{}) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Sugar().Fatalw(message, kvs...)
}
