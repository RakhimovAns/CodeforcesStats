package di

import (
	"context"
	"time"

	diut "github.com/RakhimovAns/CodeforcesStats/pkg/di"
	"github.com/RakhimovAns/CodeforcesStats/pkg/ratelimit"
)

func (d *DI) RateLimiter(ctx context.Context) *ratelimit.RateLimiter {
	return diut.Once(ctx, func(ctx context.Context) *ratelimit.RateLimiter {
		return ratelimit.NewRateLimiter(d.Redis(ctx), 10, time.Minute)
	})
}
