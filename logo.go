package logo

import (
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

type Logger struct {
	logger *zap.Logger
}

func New(level string, initialFields ...map[string]interface{}) *Logger {

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

	return &Logger{
		logger: zap.Must(config.Build()),
	}
}

func (l Logger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		l.logger.Info(msg)
	} else {
		l.logger.Info(msg, fieldsZap(fields[0])...)
	}
}

func (l Logger) Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		l.logger.Debug(msg)
	} else {
		l.logger.Debug(msg, fieldsZap(fields[0])...)
	}
}

func (l Logger) Error(msg string, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		l.logger.Error(msg)
	} else {
		l.logger.Error(msg, fieldsZap(fields[0])...)
	}
}

func (l Logger) Warning(msg string, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		l.logger.Warn(msg)
	} else {
		l.logger.Warn(msg, fieldsZap(fields[0])...)
	}
}

func (l Logger) Fatal(msg string, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		l.logger.Fatal(msg)
	} else {
		l.logger.Fatal(msg, fieldsZap(fields[0])...)
	}
}

func (l Logger) Close() {
	l.logger.Sync()
}

func fieldsZap(fields map[string]interface{}) (zapFields []zapcore.Field) {
	for k, v := range fields {
		zapFields = append(zapFields, zapcore.Field{Key: k, Interface: v})
	}

	return
}
