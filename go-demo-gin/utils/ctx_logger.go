package utils

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, l *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

func LoggerFrom(ctx context.Context) *logrus.Entry {
	if v := ctx.Value(loggerKey{}); v != nil {
		if l, ok := v.(*logrus.Entry); ok && l != nil {
			return l
		}
	}
	// fallback: dùng logger mặc định
	return logrus.NewEntry(logrus.StandardLogger())
}

// LogCtx: log theo context; nếu muốn thêm field thì truyền qua fields.
func LogCtx(ctx context.Context, level logrus.Level, message string, fields logrus.Fields) {
	LoggerFrom(ctx).WithFields(fields).Log(level, message)
}
