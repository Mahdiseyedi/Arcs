package redis

import (
	"arcs/internal/configs"
	"github.com/redis/go-redis/v9"
)

type Cli struct {
	Client *redis.Client
}

func NewRedisCli(cfg configs.Config) *Cli {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &Cli{
		Client: cli,
	}
}
