package contest

import telegramhandler "github.com/RakhimovAns/CodeforcesStats/internal/controller/handlers/telegram/deps"

type Handler struct {
	telegram telegramhandler.Telegram
	bot      telegramhandler.Bot
}

func New(tg telegramhandler.Telegram, bot telegramhandler.Bot) *Handler {
	return &Handler{
		telegram: tg,
		bot:      bot,
	}
}
