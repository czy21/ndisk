package cache

import (
	"context"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"strings"
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
		redisNodes := viper.GetString("cache.redis.cluster.nodes")
		redisUrl := viper.GetString("cache.redis.url")
		var clusterOpt redis.ClusterOptions
		var singleOpt, err = redis.ParseURL(redisUrl)
		if redisNodes != "" {
			clusterOpt = redis.ClusterOptions{
				Addrs:    strings.Split(redisNodes, ","),
				Password: viper.GetString("cache.redis.password"),
			}
			redisClient := redis.NewClusterClient(&clusterOpt)
			err = redisClient.ForEachMaster(context.Background(), func(ctx context.Context, shard *redis.Client) error {
				return shard.Ping(ctx).Err()
			})
			exception.Check(err)
			Client = Redis{ClusterClient: redisClient}
			return
		}
		exception.Check(err)
		Client = Redis{
			Client: redis.NewClient(singleOpt),
		}
	}
}
