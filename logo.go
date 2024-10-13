package logo

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMapper = map[string]zapcore.Level{
	"DEBUG":   zap.DebugLevel,
	"INFO":    zap.InfoLevel,
	"WARNING": zap.WarnLevel,
	"ERROR":   zap.ErrorLevel,
	"FATAL":   zap.FatalLevel,
}

type ctxKey struct{}

func New(level string, initialFields ...map[string]interface{}) *zap.Logger {

	levelLog, ok := levelMapper[level]
	if !ok {
		levelLog = zap.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "datetime"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(levelLog),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	if len(initialFields) > 0 {
		config.InitialFields = initialFields[0]
	}

	return zap.Must(config.Build())
}

func MapFields(fields map[string]interface{}) (zapFields []zapcore.Field) {
	for k, v := range fields {
		zapFields = append(zapFields, zapcore.Field{Key: k, Interface: v})
	}
	return
}

func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := New(""); l != nil {
		return l
	}

	return zap.NewNop()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
