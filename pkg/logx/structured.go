package logx

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Infof(ctx context.Context, format string, message ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	info(fmt.Sprintf(format, message...), traceInfo)
}

func Info(ctx context.Context, message string) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	info(message, traceInfo)
}

func Debugf(ctx context.Context, format string, message ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	debug(fmt.Sprintf(format, message...), traceInfo)
}

func Debug(ctx context.Context, message string) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	debug(message, traceInfo)
}

func Warnf(ctx context.Context, format string, message ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	warn(fmt.Sprintf(format, message...), traceInfo)
}

func Warn(ctx context.Context, message string) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	warn(message, traceInfo)
}

func Errorf(ctx context.Context, cause error, format string, message ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	err(cause, fmt.Sprintf(format, message...), traceInfo)
}

func Error(ctx context.Context, cause error, message string) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	err(cause, message, traceInfo)
}

func Fatalf(ctx context.Context, cause error, format string, message ...interface{}) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	fatal(cause, fmt.Sprintf(format, message...), traceInfo)
}

func Fatal(ctx context.Context, cause error, message string) {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)
	fatal(cause, message, traceInfo)
}

func WithContext(ctx context.Context) *zap.Logger {
	span, _ := tracer.SpanFromContext(ctx)
	traceInfo := traceInfoFromSpan(span)

	return zap.L().With(
		traceInfo.TraceID,
		traceInfo.SpanID,
		traceInfo.DatadogTraceID,
		traceInfo.DatadogSpanID)
}

func Log() *zap.SugaredLogger {
	return zap.L().Sugar()
}

func info(message string, trace DatadogTracing) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Info(message)
}

func debug(message string, trace DatadogTracing) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Debug(message)
}

func warn(message string, trace DatadogTracing) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
	).Warn(message)
}

func err(cause error, message string, trace DatadogTracing) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
		zap.String("cause", cause.Error()),
	).Error(message)
}

func fatal(cause error, message string, trace DatadogTracing) {
	zap.L().With(
		trace.TraceID,
		trace.SpanID,
		trace.DatadogTraceID,
		trace.DatadogSpanID,
		zap.String("cause", cause.Error()),
	).Fatal(message)
}
