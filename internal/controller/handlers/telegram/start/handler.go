package start

import (
	"context"

	telegramhandler "github.com/RakhimovAns/CodeforcesStats/internal/controller/handlers/telegram/deps"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

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

func (h *Handler) HandleStart(ctx *telegohandler.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	msg := update.Message

	_, err := h.telegram.SendMessage(
		context.Background(),
		msg.Chat.ID,
		"Hello, My friend. It's codeforces stats bot. Use /commands to get full list of commands.",
	)

	if err != nil {
		return err
	}

	return nil
}

//func (h *Handler) handleTestStart(ctx *telegohandler.Context, message telego.Message) error {
//	keyboard := telegoutil.InlineKeyboard(
//		telegoutil.InlineKeyboardRow(
//			telegoutil.InlineKeyboardButton("Test").
//				WithWebApp(&telego.WebAppInfo{
//					URL: config.GetTestAppURL(),
//				}),
//		),
//		telegoutil.InlineKeyboardRow(
//			telegoutil.InlineKeyboardButton(config.BotJoinCommunityText()).WithURL(communityURL),
//		),
//	)
//
//	params := telegoutil.Photo(
//		message.Chat.ChatID(),
//		telegoutil.FileFromURL(config.BotStartImageURL()),
//	)
//	params.WithCaption("Test version, hey")
//	params.WithReplyMarkup(keyboard)
//
//	_, err := ctx.Bot().SendPhoto(ctx, params)
//	if err != nil {
//		return fmt.Errorf("failed to send start message with photo: %w", err)
//	}
//
//	return nil
//}
