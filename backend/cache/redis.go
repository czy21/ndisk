package cache

import (
	"context"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	Client        *redis.Client
	ClusterClient *redis.ClusterClient
}

func (c Redis) Set(ctx context.Context, key string, value interface{}) {
	if c.ClusterClient != nil {
		_, err := c.ClusterClient.Set(ctx, key, value, time.Duration(3000)*1000).Result()
		exception.Check(err)
		return
	}
	c.Client.Set(ctx, key, value, time.Duration(3000)*10000)
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
