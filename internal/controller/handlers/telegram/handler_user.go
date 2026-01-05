package telegramhandler

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/RakhimovAns/CodeforcesStats/pkg/telegram"
	utils "github.com/RakhimovAns/CodeforcesStats/pkg/utils/user"
	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func (h *Handler) handleUser(ctx *telegohandler.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}
	msg := update.Message.Text
	parts := strings.Split(msg, " ")
	if len(parts) < 2 {
		_, err := h.telegram.SendMessage(ctx.Context(), update.Message.Chat.ID, "wrong using of this command, use /command to get more information")
		if err != nil {
			return slerr.WithSource(err, "sending message")
		}
		return nil
	}

	username := parts[len(parts)-1]

	users, err := h.bot.FetchUserInfo(ctx.Context(), username)
	if err != nil {
		slog.Error("failed to fetch user info", "user", username, "err", err)
		_, err = h.telegram.SendMessage(ctx.Context(), update.Message.Chat.ID, fmt.Sprintf("user with handle %s not found", username))
		if err != nil {
			return slerr.WithSource(err, "sending message, when fetching user info")
		}

		return nil
	}

	if len(users) == 0 {
		_, err = h.telegram.SendMessage(ctx.Context(), update.Message.Chat.ID, fmt.Sprintf("user with handle %s not found", username))
		if err != nil {
			return slerr.WithSource(err, "sending message, where user = 0 ")
		}

		return nil
	}
	user := users[0]

	text := utils.FormatUserInfo(user)
	_, err = h.telegram.SendMessage(ctx.Context(), update.Message.Chat.ID, text, telegram.WithPhoto(user.TitlePhoto))
	if err != nil {
		return slerr.WithSource(err, "sending message, when photo")
	}

	return nil
}
