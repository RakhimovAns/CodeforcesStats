package telegram

import (
	"context"
	"fmt"
	"time"

	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

func (t *Telegram) SavePreparedInlineMessage(ctx context.Context, userID int64, title, description, messageText string) (*telego.PreparedInlineMessage, error) {
	return t.savePreparedInlineMessage(ctx, userID, title, description, messageText, nil)
}

func (t *Telegram) SavePreparedInlineMessageWithButton(ctx context.Context, userID int64, title, description, messageText, buttonText, buttonURL string) (*telego.PreparedInlineMessage, error) {
	var kb *telego.InlineKeyboardMarkup
	if buttonText != "" && buttonURL != "" {
		kb = telegoutil.InlineKeyboard(
			telegoutil.InlineKeyboardRow(
				telegoutil.InlineKeyboardButton(buttonText).WithURL(buttonURL),
			),
		)
	}

	return t.savePreparedInlineMessage(ctx, userID, title, description, messageText, kb)
}

func (t *Telegram) savePreparedInlineMessage(ctx context.Context, userID int64, title, description, messageText string, replyMarkup *telego.InlineKeyboardMarkup) (*telego.PreparedInlineMessage, error) {
	params := &telego.SavePreparedInlineMessageParams{
		UserID: userID,
		Result: &telego.InlineQueryResultArticle{
			Type:  telego.ResultTypeArticle,
			ID:    fmt.Sprintf("giveaway_%d_%d", userID, time.Now().Unix()),
			Title: title,
			InputMessageContent: &telego.InputTextMessageContent{
				MessageText: messageText,
				ParseMode:   telego.ModeHTML,
			},
			Description: description,
			ReplyMarkup: replyMarkup,
		},
		AllowUserChats:    true,
		AllowBotChats:     false,
		AllowGroupChats:   true,
		AllowChannelChats: true,
	}

	out, err := t.bot.SavePreparedInlineMessage(ctx, params)
	if err != nil {
		return nil, slerr.WithSource(err)
	}

	return out, nil
}
