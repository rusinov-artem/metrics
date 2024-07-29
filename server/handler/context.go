package handler

import (
	"context"
	"time"
)

func (h *Handler) context(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}
