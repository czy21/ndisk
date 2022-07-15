package cache

import (
	"context"
	"time"
)

type Memory struct {
}

func (c Memory) Set(ctx context.Context, key string, value interface{}) {
	panic("implement me")
}

func (c Memory) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	panic("implement me")
}

func (c Memory) Get(ctx context.Context, key string) string {
	panic("implement me")
}

func (c Memory) GetObj(ctx context.Context, key string, v interface{}) {
	panic("implement me")
}

func (c Memory) Del(ctx context.Context, key string) {
	panic("implement me")
}
