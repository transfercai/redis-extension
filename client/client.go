package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var (
	once   sync.Once
	client *Client
)

type Client struct {
	*redis.Client
}

func DoInjectRedis(rc *redis.Client) *Client {
	once.Do(
		func() { client = &Client{rc} })
	return client
}

func (c *Client) GetAndDel(ctx context.Context, key string, value interface{}) (interface{}, error) {
	luaScript :=
		`if redis.call("get",KEYS[1]) == ARGV[1] then
    		return redis.call("del",KEYS[1])
		else
    		return 0
		end`
	keys := []string{key}
	return c.Eval(luaScript, keys, value).Result()
}

func (c *Client) IncrAndExpire(ctx context.Context, key string, timeSec int64) (interface{}, error) {
	luaScript :=
		`local current = redis.call("incr",KEYS[1])
		redis.call("expire",KEYS[1],ARGV[1])
		return current`
	keys := []string{key}
	return c.Eval(luaScript, keys, timeSec).Result()
}

func (c *Client) DecrAndExpire(ctx context.Context, key string, timeSec int64) (interface{}, error) {
	luaScript :=
		`local current = redis.call("decr",KEYS[1])
    	redis.call("expire",KEYS[1],ARGV[1])
		return current`
	keys := []string{key}
	return c.Eval(luaScript, keys, timeSec).Result()
}

func (c *Client) HSetAndExpire(ctx context.Context, key string, childKey string, value interface{}, timeSec int64) (interface{}, error) {
	luaScript :=
		`local ret = redis.call("hset", KEYS[1], KEYS[2], ARGV[1])
		if ret > 0 then
			return redis.call("expire", KEYS[1], ARGV[2])
		else 
			return -1
		end`
	keys := []string{key, childKey}
	return c.Eval(luaScript, keys, value, timeSec).Result()
}

func (c *Client) SetNxAndExpire(ctx context.Context, key string, value interface{}, timeSec int64) (interface{}, error) {
	luaScript :=
		`local ret = redis.call("setnx", KEYS[1], ARGV[1])
		if ret > 0 then
			return redis.call("expire", KEYS[1], ARGV[2])
		else 
			return -1
		end`
	keys := []string{key}
	return c.Eval(luaScript, keys, value, timeSec).Result()
}

func (c *Client) MSetAndExpire(ctx context.Context, keys []string, values []interface{}, timeSec int64) (interface{}, error) {
	luaScript := fmt.Sprintf(
		`local vals = ARGV local ex = %d `, timeSec)
	luaScript +=
		`for i,v in ipairs(KEYS) do
			redis.call("set",v,vals[i])
			local ret=redis.call("expire",v,ex)
			if ret ~= 1 then 
				return ret
			end
		end
		return 1
		`
	return c.Eval(luaScript, keys, values...).Result()
}

func (c *Client) GetAndExpire(ctx context.Context, key string, timeSec int64) (interface{}, error) {
	luaScript :=
		`redis.call("expire", KEYS[1], ARGV[1])
		return redis.call("get",KEYS[1])
		`
	keys := []string{key}
	return c.Eval(luaScript, keys, timeSec).Result()
}
