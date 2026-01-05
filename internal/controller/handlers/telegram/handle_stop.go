package telegramhandler

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func (h *Handler) handleStop(ctx *telegohandler.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	msg := update.Message

	_, err := h.telegram.SendMessage(
		context.Background(),
		msg.Chat.ID,
		"ОП ОП поймал стоп",
		//telego.Repl
	)
	if err != nil {
		return err
	}

	return nil
}
