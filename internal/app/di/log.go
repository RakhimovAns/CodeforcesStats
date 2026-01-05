package di

import (
	"context"
	"log/slog"

	diut "github.com/RakhimovAns/CodeforcesStats/pkg/di"
	"github.com/RakhimovAns/logger/pkg/logger/sl"
)

func (d *DI) Log(ctx context.Context) *slog.Logger {
	return diut.Once(ctx, func(ctx context.Context) *slog.Logger {
		return sl.Default()
	})
}
