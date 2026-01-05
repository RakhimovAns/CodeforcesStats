package telegram

import (
	"context"
	"fmt"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func (t *Telegram) Validate(sign string, expIn time.Duration) error {
	if sign == "" {
		return fmt.Errorf("empty sign provided")
	}
	if expIn <= 0 {
		return fmt.Errorf("invalid expiration time: %v", expIn)
	}

	return initdata.Validate(sign, t.bot.Token(), expIn)
}

func (t *Telegram) ValidateThirdPartySign(sign string, botID int64, expIn time.Duration) error {
	return initdata.ValidateThirdParty(sign, botID, expIn)
}

func (t *Telegram) FetchInitData(_ context.Context, sign string) (initdata.InitData, error) {
	if sign == "" {
		return initdata.InitData{}, fmt.Errorf("empty sign provided")
	}
	return initdata.Parse(sign)
}
