package model

type User struct {
	Handle       string  `json:"handle"`
	Email        *string `json:"email,omitempty"`
	VkID         *string `json:"vkId,omitempty"`
	OpenID       *string `json:"openId,omitempty"`
	FirstName    *string `json:"firstName,omitempty"`
	LastName     *string `json:"lastName,omitempty"`
	Country      *string `json:"country,omitempty"`
	City         *string `json:"city,omitempty"`
	Organization *string `json:"organization,omitempty"`
	Contribution int64   `json:"contribution"`
	Rank         string  `json:"rank"`
	Rating       int64   `json:"rating"`
	MaxRank      string  `json:"maxRank"`
	MaxRating    int64   `json:"maxRating"`
	Avatar       string  `json:"avatar"`
	TitlePhoto   string  `json:"titlePhoto"`
}
