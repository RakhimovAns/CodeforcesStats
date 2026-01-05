package di

import (
	"context"

	bot "github.com/RakhimovAns/CodeforcesStats/internal/service/telegram"
	diut "github.com/RakhimovAns/CodeforcesStats/pkg/di"
)

func (d *DI) BotService(ctx context.Context) *bot.Service {
	return diut.Once(ctx, func(ctx context.Context) *bot.Service {
		return bot.New(d.ExternalApiService(ctx))
	})
}
