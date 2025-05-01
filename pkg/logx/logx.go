package logx

import (
	"log"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

/*
logger for application
*/
func Init(name, version, env string) *zap.Logger {
	var cfg zap.Config

	cfg = zap.NewDevelopmentConfig()

	if env != "local" {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout"}
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.StacktraceKey = zapcore.OmitKey
		cfg.InitialFields = map[string]interface{}{
			"service-name":    name,
			"service-version": version,
		}
	}

	logger, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		log.Fatalf("build log: %s", err)
	}

	zap.ReplaceGlobals(logger)

	return logger
}

/*
logger for unit test
*/
func InitNop() *zap.Logger {
	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

	return logger
}

type DatadogTracing struct {
	TraceID        zap.Field
	SpanID         zap.Field
	DatadogTraceID zap.Field
	DatadogSpanID  zap.Field
}

func traceInfoFromSpan(span ddtrace.Span) DatadogTracing {
	var (
		spanCtx = span.Context()
		spanID  = spanCtx.SpanID()
		traceID = spanCtx.TraceID()
	)

	return DatadogTracing{
		TraceID:        zap.String("trace-id", strconv.FormatUint(traceID, 10)),
		SpanID:         zap.String("span-id", strconv.FormatUint(spanID, 10)),
		DatadogTraceID: zap.Uint64("dd.trace_id", traceID),
		DatadogSpanID:  zap.Uint64("dd.span_id", spanID),
	}
}
