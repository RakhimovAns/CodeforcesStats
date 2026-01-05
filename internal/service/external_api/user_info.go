package external_api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/RakhimovAns/CodeforcesStats/internal/model"
	"github.com/RakhimovAns/CodeforcesStats/pkg/utils"
	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/pkg/errors"
)

const (
	OkStatus = "OK"
)

func (s *Service) FetchUserInfo(ctx context.Context, username string) ([]model.User, error) {
	userLink := fmt.Sprintf("https://codeforces.com/api/user.info?handles=%s", username)
	slog.Info("fetching user info", "url", userLink)

	// Теперь получаем статус код, тело и ошибку
	statusCode, body, err := utils.Request(ctx, userLink)
	if err != nil {
		return nil, slerr.WithSource(err)
	}

	// Проверяем статус код
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", statusCode, string(body))
	}

	// Логируем тело для отладки
	slog.Debug("response body", "body", string(body))

	var rawStruct struct {
		Status string       `json:"status"`
		Result []model.User `json:"result"`
	}

	slog.Info("parsing user info")

	// Используем body ([]byte) для парсинга JSON
	if err := json.Unmarshal(body, &rawStruct); err != nil {
		slog.Error("JSON parsing failed", "error", err, "body_preview", string(body[:min(100, len(body))]))
		return nil, slerr.WithSource(errors.Wrap(err, "failed to parse JSON response"))
	}

	// Проверяем статус ответа API
	if rawStruct.Status != OkStatus {
		return nil, fmt.Errorf("API returned error status: %s", rawStruct.Status)
	}

	slog.Info("fetched user info", "users_count", len(rawStruct.Result))
	return rawStruct.Result, nil
}
