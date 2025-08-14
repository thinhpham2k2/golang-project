package utils

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type locKey struct{}

func WithLocalizer(ctx context.Context, l *i18n.Localizer) context.Context {
	return context.WithValue(ctx, locKey{}, l)
}

func LocalizerFrom(ctx context.Context) *i18n.Localizer {
	if v := ctx.Value(locKey{}); v != nil {
		if l, ok := v.(*i18n.Localizer); ok && l != nil {
			return l
		}
	}
	return nil // hoặc trả về localizer mặc định nếu bạn muốn
}
