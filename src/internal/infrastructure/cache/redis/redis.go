package redis

import (
	"context"
	domainErrors "suscord/internal/domain/errors"
	"time"

	pkgErrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	*redis.Client
}

func NewStorage(addr, password string, db int) *redisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisClient{Client: client}
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", domainErrors.ErrRedisNil
		}
		return "", pkgErrors.WithStack(err)
	}
	return value, nil
}

func (r *redisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	err := r.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (r *redisClient) Remove(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		if err == redis.Nil {
			return domainErrors.ErrRedisNil
		}
		return pkgErrors.WithStack(err)
	}
	return nil
}
