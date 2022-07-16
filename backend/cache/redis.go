package cache

import (
	"context"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

type Redis struct {
	Cmd redis.Cmdable
}

func (c Redis) Set(ctx context.Context, key string, value interface{}) {
	expiration := time.Duration(viper.GetInt64("cache.expire")) * time.Second
	_, err := c.Cmd.Set(ctx, key, value, expiration).Result()
	exception.Check(err)
}

func (c Redis) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	panic("implement me")
}

func (c Redis) Get(ctx context.Context, key string) string {
	val, err := c.Cmd.Get(ctx, key).Result()
	exception.Check(err)
	return val
}

func (c Redis) GetObj(ctx context.Context, key string, v interface{}) {
	panic("implement me")
}

func (c Redis) Del(ctx context.Context, key string) {
	panic("implement me")
}
