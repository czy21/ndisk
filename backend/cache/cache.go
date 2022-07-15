package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{})
	SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration)
	Get(ctx context.Context, key string) string
	GetObj(ctx context.Context, key string, v interface{})
	Del(ctx context.Context, key string)
}

type Client struct {
}

func (c Client) Set(ctx context.Context, key string, value interface{}) {
	panic("implement me")
}

func (c Client) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	panic("implement me")
}

func (c Client) Get(ctx context.Context, key string) string {
	panic("implement me")
}

func (c Client) GetObj(ctx context.Context, key string, v interface{}) {
	panic("implement me")
}

func (c Client) Del(ctx context.Context, key string) {
	panic("implement me")
}
