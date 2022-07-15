package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	Client *redis.Client
}

func (c Redis) Set(ctx context.Context, key string, value interface{}) {
	c.Client.Set(ctx, key, value, time.Duration(3000))
}

func (c Redis) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	panic("implement me")
}

func (c Redis) Get(ctx context.Context, key string) string {
	val, _ := c.Client.Get(ctx, key).Result()
	return val
}

func (c Redis) GetObj(ctx context.Context, key string, v interface{}) {
	panic("implement me")
}

func (c Redis) Del(ctx context.Context, key string) {
	panic("implement me")
}
