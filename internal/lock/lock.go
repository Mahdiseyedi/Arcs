package lock

import (
	redisCli "arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/utils/errmsg"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type Lock struct {
	cfg        configs.Config
	redisCli   *redisCli.Cli
	value      string
	expiration time.Duration
}

func NewLock(
	cfg configs.Config,
	cli *redisCli.Cli,
) *Lock {
	exp := time.Duration(cfg.Basic.RepublishLockDuration) * time.Second
	return &Lock{
		cfg:        cfg,
		redisCli:   cli,
		value:      uuid.NewString(),
		expiration: exp,
	}
}

func (l *Lock) AcquireLock(ctx context.Context, key string) (bool, error) {
	return l.redisCli.Client.SetNX(ctx, key, l.value, l.expiration).Result()
}

func (l *Lock) ReleaseLock(ctx context.Context, key string) error {
	val, err := l.redisCli.Client.GetDel(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	if val == l.value {
		_, err := l.redisCli.Client.Del(ctx, key).Result()
		return err
	}

	return nil
}

func (l *Lock) ExtendLock(ctx context.Context, key string) error {
	//check lock still for this instance or not
	val, err := l.redisCli.Client.Get(ctx, key).Result()
	if err != nil {
		//lock not available or expired alrdy
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	//check lock for current instance with uuid val check
	if val != l.value {
		return nil
	}

	//extend lock TTL
	ok, err := l.redisCli.Client.Expire(ctx, key, l.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errmsg.FailedExtendLock
	}
	
	return nil
}
