package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	rdb    *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(rdb *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{rdb: rdb, limit: limit, window: window}
}

func (rl *RateLimiter) Allow(ctx context.Context, userID int64) (bool, error) {
	now := time.Now().Unix()
	windowKey := now / int64(rl.window.Seconds())

	key := fmt.Sprintf("rate_limit:%d:%d", userID, windowKey)

	pipe := rl.rdb.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return incr.Val() <= int64(rl.limit), nil
}
