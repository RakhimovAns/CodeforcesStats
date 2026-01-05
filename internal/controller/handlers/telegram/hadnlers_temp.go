package telegramhandler

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/pkg/telegram"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func (h *Handler) handleTest(ctx *telegohandler.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	msg := update.Message

	_, err := h.telegram.SendMessage(
		context.Background(),
		msg.Chat.ID,
		"Выбери хей",
		telegram.WithRemoveKeyboard(),
	)

	if err != nil {
		return err
	}

	return nil
}
