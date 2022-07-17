package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/czy21/ndisk/exception"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
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

func (c Redis) SetObj(ctx context.Context, key string, value interface{}) {
	val, err := json.Marshal(value)
	exception.Check(err)
	c.Set(ctx, key, val)
}

func (c Redis) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	_, err := c.Cmd.SetEX(ctx, key, value, expiration).Result()
	exception.Check(err)
}

func (c Redis) SetObjEX(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	val, err := json.Marshal(value)
	exception.Check(err)
	c.SetEX(ctx, key, val, expiration)
}

func (c Redis) Get(ctx context.Context, key string) string {
	val, err := c.Cmd.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Tracef("%s does not exist", key)
	} else {
		exception.Check(err)
	}
	return val
}
func (c Redis) GetEX(ctx context.Context, key string, expiration time.Duration) string {
	val, err := c.Cmd.GetEx(ctx, key, expiration).Result()
	if err == redis.Nil {
		log.Tracef("%s does not exist", key)
	} else {
		exception.Check(err)
	}
	return val
}

func (c Redis) GetObj(ctx context.Context, key string, v interface{}) bool {
	val := c.Get(ctx, key)
	if val != "" {
		err := json.Unmarshal([]byte(val), v)
		exception.Check(err)
		return true
	}
	return false
}

func (c Redis) GetObjEX(ctx context.Context, key string, v interface{}, expiration time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (c Redis) Del(ctx context.Context, key string) {
	c.Cmd.Del(ctx, key)
}

func (c Redis) DelPrefix(ctx context.Context, prefix string) {
	iter := c.Cmd.Scan(ctx, 0, fmt.Sprintf("%s*", prefix), 0).Iterator()
	for iter.Next(ctx) {
		c.Del(ctx, iter.Val())
	}
}
