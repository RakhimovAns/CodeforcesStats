package telegramhandler

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/internal/model"
	"github.com/RakhimovAns/CodeforcesStats/pkg/ratelimit"
	"github.com/RakhimovAns/CodeforcesStats/pkg/telegram"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

type Telegram interface {
	SendMessage(ctx context.Context, chatID int64, text string, options ...telegram.MsgOption) (*telego.Message, error)
	CheckChannelMember(ctx context.Context, userID int64, channelUsername string) (bool, error)
}

type Bot interface {
	FetchUserInfo(ctx context.Context, username string) ([]model.User, error)
}
type Handler struct {
	telegram    Telegram
	rateLimiter ratelimit.RateLimiter
	bot         Bot
}

func New(telegram Telegram, limiter ratelimit.RateLimiter, bot Bot) *Handler {
	return &Handler{
		telegram:    telegram,
		rateLimiter: limiter,
		bot:         bot,
	}
}

func (h *Handler) Setup(th *telegohandler.BotHandler) {
	th.Handle(h.withRateLimit(h.handleStart), telegohandler.CommandEqual("start"))
	th.Handle(h.withRateLimit(h.handleTest), telegohandler.CommandEqual("test"))
	th.Handle(h.withRateLimit(h.handleUser), telegohandler.CommandEqual("user"))
	//th.Handle(h.withRateLimit(h.handleUnknownCommand), telegohandler.Command())

	th.Handle(h.withRateLimit(h.handleStop), telegohandler.TextEqual("СТОП"))

	// любой текст
	th.Handle(h.withRateLimit(h.handleText))
	//th.HandleMessage(h.handleTestStart, telegohandler.CommandEqual("starttest"))
	//th.HandleMessage(h.createNotificationsChat, telegohandler.CommandEqual("add_insufficient_balance_chat"))
}

//func allowedChatTypePredicate() telegohandler.Predicate {
//	return func(ctx context.Context, update telego.Update) bool {
//		return slices.Contains([]string{telego.ChatTypeSupergroup, telego.ChatTypeChannel}, update.MyChatMember.Chat.Type)
//	}
//}
//
//func nonEmptyChatUsernamePredicate() telegohandler.Predicate {
//	return func(ctx context.Context, update telego.Update) bool {
//		return update.MyChatMember.Chat.Username != ""
//	}
//}
