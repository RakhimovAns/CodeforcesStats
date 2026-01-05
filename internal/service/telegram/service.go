package bot

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/internal/model"
)

type UserAPI interface {
	FetchUserInfo(ctx context.Context, username string) ([]model.User, error)
}

type Service struct {
	userAPI UserAPI
}

func New(userAPI UserAPI) *Service {
	return &Service{
		userAPI: userAPI,
	}
}
