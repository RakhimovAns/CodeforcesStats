package external_api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/RakhimovAns/CodeforcesStats/internal/model"
	model2 "github.com/RakhimovAns/CodeforcesStats/internal/service/model"
	"github.com/RakhimovAns/CodeforcesStats/pkg/utils"
	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/pkg/errors"
)

func (s *Service) FetchContestInfo(ctx context.Context) ([]model.Contest, error) {
	link := fmt.Sprintf("https://codeforces.com/api/contest.list?gym=false")
	slog.Info("fetching contest info", "url", link)
	statusCode, body, err := utils.Request(ctx, link)
	if err != nil {
		return nil, slerr.WithSource(err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", statusCode, string(body))
	}

	slog.Debug("response body", "body", string(body))

	var rawStruct struct {
		Status string          `json:"status"`
		Result []model.Contest `json:"result"`
	}

	if err := json.Unmarshal(body, &rawStruct); err != nil {
		slog.Error("JSON parsing failed", "error", err, "body_preview", string(body[:min(100, len(body))]))
		return nil, slerr.WithSource(errors.Wrap(err, "failed to parse JSON response"))
	}

	if rawStruct.Status != model2.OkStatus {
		return nil, fmt.Errorf("API returned error status: %s", rawStruct.Status)
	}

	slog.Info("fetched contest info", "contests_count", len(rawStruct.Result))
	return rawStruct.Result, nil
}
