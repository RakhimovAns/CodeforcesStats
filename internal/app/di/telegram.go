package di

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/internal/config"
	telegramhandler "github.com/RakhimovAns/CodeforcesStats/internal/controller/handlers/telegram"
	diut "github.com/RakhimovAns/CodeforcesStats/pkg/di"
	"github.com/RakhimovAns/CodeforcesStats/pkg/telegram"
)

func (d *DI) Telegram(ctx context.Context) *telegram.Telegram {
	return diut.Once(ctx, func(ctx context.Context) *telegram.Telegram {
		tg, err := telegram.New(
			config.Telegram(),
		)
		if err != nil {
			d.mustExit(err)
		}

		return tg
	})
}

func (d *DI) TelegramHandler(ctx context.Context) *telegramhandler.Handler {
	return diut.Once(ctx, func(ctx context.Context) *telegramhandler.Handler {
		return telegramhandler.New(
			d.Telegram(ctx),
			*d.RateLimiter(ctx),
			d.ExternalApiService(ctx),
		)
	})
}

func (d *DI) Notifier(ctx context.Context) *telegram.Telegram {
	return diut.Once(ctx, func(ctx context.Context) *telegram.Telegram {
		return d.Telegram(ctx)
	})
}

//func (d *DI) NotificationBot(ctx context.Context) *nfbot.NotificationBot {
//	return diut.Once(ctx, func(ctx context.Context) *nfbot.NotificationBot {
//		token := config.TelegramNotificationBotToken()
//		channelID := config.TelegramNotificationBotChannelID()
//
//		mainLogger := d.Log(ctx)
//		notificationLogger := mainLogger.With(slog.String("component", "notification_bot"))
//
//		bot, err := nfbot.New(nfbot.Config{
//			Token:     token,
//			ChannelID: channelID,
//			Logger:    notificationLogger,
//		})
//		if err != nil {
//			return nil
//		}
//
//		return bot
//	})
//}
