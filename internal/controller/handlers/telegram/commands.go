package telegramhandler

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

// обработчик для неизвестных команд
//func (h *Handler) handleUnknownCommand(ctx *telegohandler.Context, update telego.Update) error {
//	if update.Message == nil {
//		return nil
//	}
//
//	_, err := h.telegram.SendMessage(
//		context.Background(),
//		update.Message.Chat.ID,
//		"❓ Неизвестная команда. Попробуйте /start",
//	)
//	return err
//}

// обработчик для произвольного текста
func (h *Handler) handleText(ctx *telegohandler.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	_, err := h.telegram.SendMessage(
		context.Background(),
		update.Message.Chat.ID,
		"Я получил текст: "+update.Message.Text,
	)
	return err
}
