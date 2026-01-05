package utils

import (
	"fmt"
	"html"

	"github.com/RakhimovAns/CodeforcesStats/internal/model"
)

func FormatUserInfo(user model.User) string {
	var result string
	result += fmt.Sprintf("<b>Handler</b> : <i>%s</i>", user.Handle)
	if user.Organization != nil {
		result += fmt.Sprintf("\n<b>Organization</b>: <i>%s</i>", *user.Organization)
	}
	result += fmt.Sprintf("\n<b>Contribution</b>: <i>%v</i>", user.Contribution)
	result += fmt.Sprintf("\n<b>Rank: %s</b>, <b>Max_Rank</b>: <s>%s</s>", user.Rank, user.MaxRank)
	result += fmt.Sprintf("\n<b>Rating: %v</b>, <b>Max_Rating</b>: <s>%v</s>", user.Rating, user.MaxRating)
	return result
}

func EscapeHTML(text string) string {
	return html.EscapeString(text)
}
