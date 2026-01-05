package di

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/internal/service/external_api"
	diut "github.com/RakhimovAns/CodeforcesStats/pkg/di"
)

func (d *DI) ExternalApiService(ctx context.Context) *external_api.Service {
	return diut.Once(ctx, func(ctx context.Context) *external_api.Service {
		return external_api.New()
	})
}
