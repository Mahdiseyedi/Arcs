package health

import (
	"arcs/internal/clients/nats/producer"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Svc struct {
	gDB      *gorm.DB
	redisCli *redis.Client
	producer *producer.Producer
}

func NewHealthSvc(
	db *gorm.DB,
	redisCli *redis.Client,
	producer *producer.Producer,
) *Svc {
	return &Svc{
		gDB:      db,
		redisCli: redisCli,
		producer: producer,
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
	return s.producer.HealthCheck()
}
