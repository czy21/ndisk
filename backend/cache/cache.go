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
		url := viper.GetString("cache.redis.url")
		var singleOpt, err = redis.ParseURL(url)
		cmd := redis.NewClient(singleOpt)
		exception.Check(err)
		Client = Redis{Cmd: cmd}
	}
}
