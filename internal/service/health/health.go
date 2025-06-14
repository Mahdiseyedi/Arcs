package health

import (
	natsCli "arcs/internal/clients/nats"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Svc struct {
	gDB        *gorm.DB
	redisCli   *redis.Client
	natsClient *natsCli.Client
}

func NewHealthSvc(
	db *gorm.DB,
	redisCli *redis.Client,
	natsClient *natsCli.Client,
) *Svc {
	return &Svc{
		gDB:        db,
		redisCli:   redisCli,
		natsClient: natsClient,
	}
}

func (s *Svc) DBHealthCheck(ctx context.Context) error {
	insideDB, err := s.gDB.DB()
	if err != nil {
		return err
	}

	return insideDB.PingContext(ctx)
}

func (s *Svc) RedisHealthCheck(ctx context.Context) error {
	return s.redisCli.Ping(ctx).Err()
}

func (s *Svc) NatsHealthCheck(ctx context.Context) error {
	return s.natsClient.HealthCheck()
}
