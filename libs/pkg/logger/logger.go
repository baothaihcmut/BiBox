package logger

import (
	"context"

	"github.com/baothaihcmut/BiBox/libs/pkg/constant"
	"github.com/sirupsen/logrus"
)

// Logger interface to support different log levels
type Logger interface {
	Debug(ctx context.Context, args map[string]any, msg string)
	Debugf(ctx context.Context, args map[string]any, template string, format ...interface{})
	Info(ctx context.Context, args map[string]any, msg string)
	Infof(ctx context.Context, args map[string]any, template string, format ...interface{})
	Warn(ctx context.Context, args map[string]any, msg string)
	Warnf(ctx context.Context, args map[string]any, template string, format ...interface{})
	Error(ctx context.Context, args map[string]any, msg string)
	Errorf(ctx context.Context, args map[string]any, template string, format ...interface{})
	Fatal(ctx context.Context, args map[string]any, msg string)
	Fatalf(ctx context.Context, args map[string]any, template string, format ...interface{})
	Panic(ctx context.Context, args map[string]any, msg string)
	Panicf(ctx context.Context, args map[string]any, template string, format ...interface{})
}

// logger is the concrete implementation of the Logger interface
type logger struct {
	l *logrus.Logger
}

// NewLogger initializes and returns a logger instance
func NewLogger(l *logrus.Logger) Logger {
	return &logger{l: l}
}

// Debug level logging with context
func (l *logger) Debug(ctx context.Context, args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value(constant.RequesIdContext); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Debug(msg)
}

func (l *logger) Debugf(ctx context.Context, args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Debugf(template, format...)
}

// Info level logging with context
func (l *logger) Info(ctx context.Context, args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Info(msg)
}

func (l *logger) Infof(ctx context.Context, args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Infof(template, format...)
}

// Warn level logging with context
func (l *logger) Warn(ctx context.Context, args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Warn(msg)
}

func (l *logger) Warnf(ctx context.Context, args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Warnf(template, format...)
}

// Error level logging with context
func (l *logger) Error(ctx context.Context, args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Error(msg)
}

func (l *logger) Errorf(ctx context.Context, args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Errorf(template, format...)
}

// Fatal level logging with context
func (l *logger) Fatal(ctx context.Context, args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Fatal(msg)
}

func (l *logger) Fatalf(ctx context.Context, args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Fatalf(template, format...)
}

// Panic level logging with context
func (l *logger) Panic(ctx context.Context, args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Panic(msg)
}

func (l *logger) Panicf(ctx context.Context, args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		args["request_id"] = requestID
	}
	l.l.WithFields(logrus.Fields(args)).Panicf(template, format...)
}
