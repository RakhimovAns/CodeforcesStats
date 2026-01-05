package telegram

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/bytedance/sonic"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoapi"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/shopspring/decimal"
)

type BotUrls struct {
	BotWebAppURL      string
	EnvironmentAppURL string
}

type BotTexts struct {
	BotOpenPortalsText string
}

type Telegram struct {
	bot *telego.Bot
}

func New(token string) (*Telegram, error) {
	bot, err := telego.NewBot(token, telego.WithAPICaller(&telegoapi.RetryCaller{
		Caller:       telegoapi.DefaultFastHTTPCaller,
		MaxAttempts:  3,
		ExponentBase: 2,
		StartDelay:   time.Millisecond * 10,
		MaxDelay:     time.Second,
		RateLimit:    telegoapi.RetryRateLimitWait,
	}))
	if err != nil {
		return nil, slerr.WithSource(err)
	}

	telego.SetJSONMarshal(sonic.Marshal)
	telego.SetJSONUnmarshal(sonic.Unmarshal)

	return &Telegram{
		bot: bot,
	}, nil
}

func (t *Telegram) Bot() *telego.Bot {
	return t.bot
}

func (t *Telegram) Ping(ctx context.Context) error {
	_, err := t.bot.GetMe(ctx)
	if err != nil {
		return slerr.WithSource(err)
	}

	return nil
}

func (t *Telegram) CheckChannelMember(ctx context.Context, userID int64, channelUsername string) (bool, error) {
	res, err := t.bot.GetChatMember(ctx, &telego.GetChatMemberParams{
		ChatID: telegoutil.Username(t.fixChannelUsername(channelUsername)),
		UserID: userID,
	})
	if err != nil {
		return false, slerr.WithSource(err)
	}

	memberStatuses := []string{
		telego.MemberStatusCreator,
		telego.MemberStatusAdministrator,
		telego.MemberStatusMember,
		telego.MemberStatusRestricted,
	}

	return slices.Contains(memberStatuses, res.MemberStatus()), nil
}

func (t *Telegram) CheckBotAdminStatus(ctx context.Context, channelUsername string) (bool, error) {
	res, err := t.bot.GetChatMember(ctx, &telego.GetChatMemberParams{
		ChatID: telegoutil.Username(t.fixChannelUsername(channelUsername)),
		UserID: t.bot.ID(),
	})
	if err != nil {
		return false, slerr.WithSource(err)
	}

	adminStatuses := []string{telego.MemberStatusCreator, telego.MemberStatusAdministrator}

	return slices.Contains(adminStatuses, res.MemberStatus()), nil
}

func (t *Telegram) GetChannelTitle(ctx context.Context, channelUsername string) (string, error) {
	res, err := t.bot.GetChat(ctx, &telego.GetChatParams{
		ChatID: telegoutil.Username(t.fixChannelUsername(channelUsername)),
	})
	if err != nil {
		return "", slerr.WithSource(err)
	}

	return res.Title, nil
}

func (t *Telegram) SendPhotoMessage(ctx context.Context, userID int64, photoURL, text string, args ...any) error {
	formattedText := fmt.Sprintf(text, args...)

	_, err := t.SendMessage(ctx, userID, formattedText, WithPhoto(photoURL))

	return err
}

//func (t *Telegram) SendMessageWithButton(ctx context.Context, userID int64, text string, args ...any) error {
//	formattedText := fmt.Sprintf(text, args...)
//
//	_, err := t.SendMessage(ctx, userID, formattedText, WithWebAppButton(t.urls, t.texts, "", false))
//
//	return err
//}
//
//func (t *Telegram) SendMessageWithEnvironmentButton(ctx context.Context, userID int64, text string, args ...any) error {
//	formattedText := fmt.Sprintf(text, args...)
//
//	_, err := t.SendMessage(ctx, userID, formattedText, WithWebAppButton(t.urls, t.texts, "", true))
//
//	return err
//}

//func (t *Telegram) SendMessageWithDeepLinkButton(ctx context.Context, userID int64, text string, route model.StartParam, buttonText string, args ...any) error {
//	formattedText := fmt.Sprintf(text, args...)
//
//	if buttonText == "" {
//		buttonText = t.texts.BotOpenPortalsText
//	}
//
//	button := DeepLinkButton{
//		Text:  buttonText,
//		Route: route,
//	}
//
//	_, err := t.SendMessage(ctx, userID, formattedText, WithDeepLinkButtons(t.urls, t.texts, button))
//
//	if err != nil {
//		return slerr.WithSource(err)
//	}
//
//	return err
//}
//
//func (t *Telegram) SendMessageWithDeepLinkButtons(ctx context.Context, userID int64, text string, buttons []DeepLinkButton, args ...any) error {
//	formattedText := fmt.Sprintf(text, args...)
//
//	_, err := t.SendMessage(ctx, userID, formattedText, WithDeepLinkButtons(t.urls, t.texts, buttons...))
//
//	return err
//}

func (t *Telegram) SendMessageToChannel(ctx context.Context, channelID int64, text string, args ...any) error {
	formattedText := fmt.Sprintf(text, args...)

	params := telegoutil.Message(telegoutil.ID(channelID), formattedText)
	params.WithParseMode(telego.ModeHTML)

	_, err := t.bot.SendMessage(ctx, params)
	if err != nil {
		return slerr.WithSource(err)
	}

	return nil
}

func (t *Telegram) SendMessageToChannelWithButton(ctx context.Context, channelID int64, text string, buttonText string, buttonURL string, args ...any) error {
	formattedText := fmt.Sprintf(text, args...)

	keyboard := telegoutil.InlineKeyboard(
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton(buttonText).WithURL(buttonURL),
		),
	)

	params := telegoutil.Message(telegoutil.ID(channelID), formattedText)
	params.WithReplyMarkup(keyboard)
	params.WithParseMode(telego.ModeHTML)

	_, err := t.bot.SendMessage(ctx, params)
	if err != nil {
		return slerr.WithSource(err)
	}

	return nil
}

func (t *Telegram) FetchStarBalance(ctx context.Context, businessConnectionID string) (decimal.Decimal, error) {
	balance, err := t.bot.GetBusinessAccountStarBalance(ctx, &telego.GetBusinessAccountStarBalanceParams{
		BusinessConnectionID: businessConnectionID,
	})
	if err != nil {
		return decimal.Decimal{}, slerr.WithSource(err)
	}

	return decimal.NewFromInt(int64(balance.Amount)), nil
}

func (t *Telegram) CheckUserAdminStatus(ctx context.Context, userID int64, channelUsername string) (bool, error) {
	res, err := t.bot.GetChatMember(ctx, &telego.GetChatMemberParams{
		ChatID: telegoutil.Username(t.fixChannelUsername(channelUsername)),
		UserID: userID,
	})
	if err != nil {
		return false, slerr.WithSource(err)
	}

	adminStatuses := []string{telego.MemberStatusCreator, telego.MemberStatusAdministrator}

	return slices.Contains(adminStatuses, res.MemberStatus()), nil
}

func (t *Telegram) GetChannelID(ctx context.Context, channelUsername string) (int64, error) {
	res, err := t.bot.GetChat(ctx, &telego.GetChatParams{
		ChatID: telegoutil.Username(t.fixChannelUsername(channelUsername)),
	})
	if err != nil {
		return 0, slerr.WithSource(err)
	}

	return res.ID, nil
}

func (t *Telegram) GetChannelUsernameByID(ctx context.Context, channelID int64) (string, error) {
	res, err := t.bot.GetChat(ctx, &telego.GetChatParams{
		ChatID: telegoutil.ID(channelID),
	})
	if err != nil {
		return "", slerr.WithSource(err)
	}

	return res.Username, nil
}

func (t *Telegram) GetUserChatBoosts(ctx context.Context, userID int64, channelUsername string) (int, error) {
	res, err := t.bot.GetUserChatBoosts(ctx, &telego.GetUserChatBoostsParams{
		ChatID: telegoutil.Username(t.fixChannelUsername(channelUsername)),
		UserID: userID,
	})
	if err != nil {
		var tgErr *telegoapi.Error
		if !errors.As(err, &tgErr) {
			return 0, slerr.WithSource(err)
		}

		if strings.Contains(tgErr.Description, "user not found") || strings.Contains(tgErr.Description, "PARTICIPANT_ID_INVALID") || strings.Contains(tgErr.Description, "USER_ID_INVALID") {
			return 0, nil
		}

		return 0, slerr.WithSource(err)
	}

	currentTime := time.Now().Unix()
	activeBoosts := 0

	for _, boost := range res.Boosts {
		isActive := boost.ExpirationDate > currentTime

		if isActive {
			activeBoosts++
		}
	}

	return activeBoosts, nil
}

func (t *Telegram) SendVideoURLToChannel(ctx context.Context, channelID int64, videoURL string, text string, args ...any) error {
	formattedText := fmt.Sprintf(text, args...)

	inputFile := telegoutil.FileFromURL(videoURL)
	params := telegoutil.Video(telegoutil.ID(channelID), inputFile)
	params.WithCaption(formattedText)
	params.WithParseMode(telego.ModeHTML)

	_, err := t.bot.SendVideo(ctx, params)
	if err != nil {
		return slerr.WithSource(err)
	}

	return nil
}

func (t *Telegram) SendVideoURLToChannelWithButton(ctx context.Context, channelID int64, videoURL string, text string, buttonText string, buttonURL string, args ...any) error {
	formattedText := fmt.Sprintf(text, args...)

	keyboard := telegoutil.InlineKeyboard(
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton(buttonText).WithURL(buttonURL),
		),
	)

	inputFile := telegoutil.FileFromURL(videoURL)
	params := telegoutil.Video(telegoutil.ID(channelID), inputFile)
	params.WithCaption(formattedText)
	params.WithParseMode(telego.ModeHTML)
	params.WithReplyMarkup(keyboard)

	_, err := t.bot.SendVideo(ctx, params)
	if err != nil {
		return slerr.WithSource(err)
	}

	return nil
}

func (t *Telegram) TransferGift(ctx context.Context, giftID string, toUserID int64, businessConnectionID string) error {
	paymentAmount := 25

	params := &telego.TransferGiftParams{
		BusinessConnectionID: businessConnectionID,
		OwnedGiftID:          giftID,
		NewOwnerChatID:       toUserID,
		StarCount:            paymentAmount,
	}

	err := t.bot.TransferGift(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to transfer gift: %w", err)
	}

	return nil
}

func (t *Telegram) fixChannelUsername(channelUsername string) string {
	if !strings.HasPrefix(channelUsername, "@") {
		channelUsername = fmt.Sprintf("@%s", channelUsername)
	}

	return channelUsername
}
