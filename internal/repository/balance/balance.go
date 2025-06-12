package balance

import (
	"arcs/internal/clients/redis"
	"context"
	"fmt"
)

type Repository struct {
	redis *redis.Cli
}

const (
	balanceKeyPrefix = "balance:"
)

func NewBalanceRepository(
	redisClient *redis.Cli,
) *Repository {
	return &Repository{
		redis: redisClient,
	}
}

func (r *Repository) balanceKey(userID string) string {
	return fmt.Sprintf("%s%s", balanceKeyPrefix, userID)
}

func (r *Repository) Set(ctx context.Context, userID string, amount int64) error {
	return r.redis.Client.Set(ctx, r.balanceKey(userID), amount, 0).Err()
}

func (r *Repository) Get(ctx context.Context, userID string) (int64, error) {
	return r.redis.Client.Get(ctx, r.balanceKey(userID)).Int64()
}

func (r *Repository) Increase(ctx context.Context, userID string, amount int64) error {
	return r.redis.Client.IncrBy(ctx, r.balanceKey(userID), amount).Err()
}

func (r *Repository) Decrease(ctx context.Context, userID string, amount int64) error {
	return r.redis.Client.DecrBy(ctx, r.balanceKey(userID), amount).Err()
}
