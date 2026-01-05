package telegramhandler

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func (h *Handler) withRateLimit(
	handler telegohandler.Handler,
) telegohandler.Handler {
	return func(ctx *telegohandler.Context, update telego.Update) error {
		if update.Message == nil {
			return handler(ctx, update)
		}

		userID := update.Message.From.ID

		allowed, err := h.rateLimiter.Allow(context.Background(), userID)
		if err != nil {
			slog.Error("rate limiter error", "error", err)
			return handler(ctx, update) // пропускаем при ошибке
		}

		if !allowed {
			_, err := h.telegram.SendMessage(
				context.Background(),
				update.Message.Chat.ID,
				"⚠️ Слишком много запросов. Пожалуйста, подождите немного.",
			)
			if err != nil {
				slog.Error("failed to send rate limit message", "error", err)
			}
			return nil
		}

		return handler(ctx, update)
	}
}
