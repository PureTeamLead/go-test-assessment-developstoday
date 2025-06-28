package logger

import (
	"context"
	"go.uber.org/zap"
	"log"
)

const (
	Key      = "logger"
	HTTPPort = "http_port"
	DevEnv   = "dev"
	ProdEnv  = "prod"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context, env string) context.Context {
	var logger *zap.Logger
	var err error

	switch env {
	case DevEnv:
		logger, err = zap.NewDevelopment()
	case ProdEnv:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("failed setting up logger: %v", err)
	}

	logg := &Logger{l: logger}

	ctx = context.WithValue(ctx, Key, logg)
	return ctx
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l *Logger) WithPort(ctx context.Context, portKey string) {
	httpPort := ctx.Value(portKey)

	if httpPort != nil {
		l.l = l.l.With(zap.Int(HTTPPort, httpPort.(int)))
		return
	}

	GetLoggerFromCtx(ctx).Info("No port specified for debugging")
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Error(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	l.l.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.l.Debug(msg, fields...)
}
