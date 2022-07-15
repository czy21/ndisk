package cache

import (
	"context"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{})
	SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration)
	Get(ctx context.Context, key string) string
	GetObj(ctx context.Context, key string, v interface{})
	Del(ctx context.Context, key string)
}

var Client Cache

func Boot() {
	cacheType := viper.GetString("cache.type")
	if cacheType == "redis" {
		opt, err := redis.ParseURL(viper.GetString("cache.redis.url"))
		exception.Check(err)
		Client = Redis{
			Client: redis.NewClient(opt),
		}
	}
}
