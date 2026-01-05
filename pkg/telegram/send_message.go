package telegram

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/internal/model"
	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

type DeepLinkButton struct {
	Text  string
	Route model.StartParam
}

type MessageOptions struct {
	PhotoURL    string
	ParseMode   string
	ReplyMarkup telego.ReplyMarkup
}

type MsgOption func(*MessageOptions)

func WithPhoto(photoURL string) MsgOption {
	return func(opts *MessageOptions) {
		opts.PhotoURL = photoURL
	}
}

func WithWebAppButton(provider BotUrls, texts BotTexts, text string, useEnvironmentURL bool) MsgOption {
	return func(opts *MessageOptions) {
		url := provider.BotWebAppURL
		if useEnvironmentURL {
			url = provider.EnvironmentAppURL
		}

		if text == "" {
			text = texts.BotOpenPortalsText
		}

		opts.ReplyMarkup = telegoutil.InlineKeyboard(
			telegoutil.InlineKeyboardRow(
				telegoutil.InlineKeyboardButton(text).WithWebApp(&telego.WebAppInfo{
					URL: url,
				}),
			),
		)
	}
}

func WithURLButtons(buttons ...map[string]string) MsgOption {
	return func(opts *MessageOptions) {
		if len(buttons) == 0 {
			return
		}

		rows := make([]telego.InlineKeyboardButton, 0, len(buttons))

		for _, button := range buttons {
			for text, url := range button {
				rows = append(rows, telegoutil.InlineKeyboardButton(text).WithURL(url))
			}
		}

		keyboard := telegoutil.InlineKeyboard(rows)

		opts.ReplyMarkup = keyboard
	}
}

//func WithDeepLinkButtons(urls BotUrls, texts BotTexts, buttons ...DeepLinkButton) MsgOption {
//	return func(opts *MessageOptions) {
//		if len(buttons) == 0 {
//			return
//		}
//
//		rows := make([][]telego.InlineKeyboardButton, 0, 3)
//		for _, btn := range buttons {
//			buttonText := btn.Text
//			if buttonText == "" {
//				buttonText = texts.BotOpenPortalsText
//			}
//
//			rows = append(rows, telegoutil.InlineKeyboardRow(
//				telegoutil.InlineKeyboardButton(buttonText).WithURL(config.GetBotDeepLink(btn.Route)),
//			))
//		}
//
//		keyboard := telegoutil.InlineKeyboard(rows...)
//
//		opts.ReplyMarkup = keyboard
//	}
//}

func WithInlineKeyboard(markup telego.ReplyMarkup) MsgOption {
	return func(opts *MessageOptions) {
		opts.ReplyMarkup = markup
	}
}

func WithReplyKeyboard(rows ...[]telego.KeyboardButton) MsgOption {
	return func(opts *MessageOptions) {
		opts.ReplyMarkup = &telego.ReplyKeyboardMarkup{
			Keyboard:        rows,
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	}
}

func WithReplyMarkup(markup telego.ReplyMarkup) MsgOption {
	return func(opts *MessageOptions) {
		opts.ReplyMarkup = markup
	}
}

func WithRemoveKeyboard() MsgOption {
	return func(opts *MessageOptions) {
		opts.ReplyMarkup = &telego.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      false,
		}
	}
}

func (t *Telegram) SendMessage(ctx context.Context, chatID int64, text string, options ...MsgOption) (*telego.Message, error) {
	opts := &MessageOptions{
		ParseMode: telego.ModeHTML,
	}

	for _, opt := range options {
		opt(opts)
	}

	var msg *telego.Message
	var err error

	if opts.PhotoURL != "" {
		photoOpts := telegoutil.Photo(
			telegoutil.ID(chatID),
			telegoutil.FileFromURL(opts.PhotoURL),
		)
		photoOpts.WithParseMode(opts.ParseMode)

		if opts.ReplyMarkup != nil {
			photoOpts.WithReplyMarkup(opts.ReplyMarkup)
		}

		if text != "" {
			photoOpts.Caption = text
		}

		msg, err = t.bot.SendPhoto(
			ctx,
			photoOpts,
		)
	} else {
		sendOpts := telegoutil.Message(telegoutil.ID(chatID), text)
		sendOpts.WithParseMode(opts.ParseMode)
		sendOpts.WithReplyMarkup(opts.ReplyMarkup)

		if opts.ReplyMarkup != nil {
			sendOpts.WithReplyMarkup(opts.ReplyMarkup)
		}

		msg, err = t.bot.SendMessage(
			ctx,
			sendOpts,
		)
	}

	if err != nil {
		return nil, slerr.WithSource(err)
	}

	return msg, nil
}
