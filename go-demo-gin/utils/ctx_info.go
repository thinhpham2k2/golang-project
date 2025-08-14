package utils

import (
	"context"
	"go-demo-gin/models"
)

type infoKey struct{}

func WithInformation(ctx context.Context, i *models.User) context.Context {
	return context.WithValue(ctx, infoKey{}, i)
}

func InformationFrom(ctx context.Context) *models.User {
	if v := ctx.Value(infoKey{}); v != nil {
		if i, ok := v.(*models.User); ok && i != nil {
			return i
		}
	}
	return nil // hoặc trả về localizer mặc định nếu bạn muốn
}
