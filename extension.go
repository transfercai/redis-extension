package redis_extension

import (
	"context"
)

type Extension interface {
	GetAndDel(ctx context.Context, key string, value interface{}) (interface{}, error)
	SetNxAndExpire(ctx context.Context, key string, value interface{}, timeSec int64) (interface{}, error)
	IncrAndExpire(ctx context.Context, key string, timeSec int64) (interface{}, error)
	DecrAndExpire(ctx context.Context, key string, timeSec int64) (interface{}, error)
	HSetAndExpire(ctx context.Context, key string, childKey string, value interface{}, timeSec int64) (interface{}, error)
	MSetAndExpire(ctx context.Context, keys []string, values []interface{}, timeSec int64) (interface{}, error)
}
