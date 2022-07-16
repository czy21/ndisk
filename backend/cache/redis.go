package cache

import (
	"context"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

type Redis struct {
	Client        *redis.Client
	ClusterClient *redis.ClusterClient
}

func (c Redis) Set(ctx context.Context, key string, value interface{}) {
	expiration := time.Duration(viper.GetInt64("cache.expire")) * time.Second
	if c.ClusterClient != nil {
		_, err := c.ClusterClient.Set(ctx, key, value, expiration).Result()
		exception.Check(err)
		return
	}
	c.Client.Set(ctx, key, value, expiration)
}

func (c Redis) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	panic("implement me")
}

func (c Redis) Get(ctx context.Context, key string) string {
	if c.ClusterClient != nil {
		val, _ := c.ClusterClient.Get(ctx, key).Result()
		return val
	}
	val, _ := c.Client.Get(ctx, key).Result()
	return val
}

func (c Redis) GetObj(ctx context.Context, key string, v interface{}) {
	panic("implement me")
}

func (c Redis) Del(ctx context.Context, key string) {
	panic("implement me")
}
