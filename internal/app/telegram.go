package app

import (
	"context"
	"log/slog"

	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/RakhimovAns/wrapper/pkg/closer"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func (a *App) runTelegramBot(ctx context.Context) error {
	log := a.di.Log(ctx).With(
		slog.String("component", "telegram_bot"),
	)

	bot := a.di.Telegram(ctx).Bot()

	if err := bot.DeleteWebhook(ctx, nil); err != nil {
		return slerr.WithSource(err)
	}

	allowedUpdates := []string{
		telego.MessageUpdates,
		telego.BusinessMessageUpdates,
		telego.MyChatMemberUpdates,
	}

	var updates <-chan telego.Update

	upds, err := bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		AllowedUpdates: allowedUpdates,
	})
	if err != nil {
		return slerr.WithSource(err)
	}

	updates = upds

	bh, err := telegohandler.NewBotHandler(bot, updates)
	if err != nil {
		return slerr.WithSource(err)
	}

	closer.Add(func() error {
		log.InfoContext(ctx, "shutting down telegram bot")

		return bh.Stop()
	})

	a.setupTelegramHandlers(ctx, bh)

	log.Info("go telegram bot!")

	if err := bh.Start(); err != nil {
		return slerr.WithSource(err)
	}

	return nil
}

func (a *App) setupTelegramHandlers(ctx context.Context, th *telegohandler.BotHandler) {
	telegramHandler := a.di.TelegramHandler(ctx)
	telegramHandler.Setup(th)
}
