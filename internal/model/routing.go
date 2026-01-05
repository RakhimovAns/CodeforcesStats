package model

type StartParam string

const (
	StartParamOffer        StartParam = "offer"
	StartParamWallet       StartParam = "wallet"
	StartParamMyGifts      StartParam = "my-gifts"
	StartParamGiveaway     StartParam = "giveaway"
	StartParamGiftDelivery StartParam = "gift-delivery"
	StartParamPayment      StartParam = "payment"
)

type DeepLinkButton struct {
	Text  string
	Route StartParam
}

func IsValidStartParam(param string) bool {
	switch StartParam(param) {
	case StartParamOffer, StartParamWallet, StartParamMyGifts, StartParamGiveaway:
		return true
	default:
		return false
	}
}
