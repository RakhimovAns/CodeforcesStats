package telegramhandler

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/internal/model"
	"github.com/RakhimovAns/CodeforcesStats/pkg/telegram"
	"github.com/mymmrac/telego"
)

type Telegram interface {
	SendMessage(ctx context.Context, chatID int64, text string, opts ...telegram.MsgOption) (*telego.Message, error)
}

type Bot interface {
	FetchUserInfo(ctx context.Context, username string) ([]model.User, error)
}
