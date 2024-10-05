package handler

import (
	"context"
	"time"
)

// context Используется для создания контекста во всех хендлерах
func (h *Handler) context(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}
