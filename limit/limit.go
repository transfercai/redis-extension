package limit

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"github.com/transfercai/redis-extension/client"
)

var (
	limitManager sync.Map
)

type Limit struct {
	Key string
	L   Limiter
}

type Limiter interface {
	Add(value string, timeStamp int64) (int64, error)
	IsLimit(start, end int64) (bool, error)
	Del(start int64) (int64, error)
	CalcStart(current int64) int64 //nano
}

type RedisLimit struct {
	rc       *client.Client
	Key      string `toml:"key"`
	Count    int64  `toml:"count"`
	Duration int64  `toml:"duration"` //default:ms
}

type HandleRDL func(rdl *RedisLimit)

func SetCount(count int64) HandleRDL {
	return func(rdl *RedisLimit) {
		rdl.Count = count
	}
}

func SetDuration(duration int64) HandleRDL {
	return func(rdl *RedisLimit) {
		rdl.Duration = duration
	}
}

func NewRedisLimit(name string, cli *client.Client, hdl ...HandleRDL) *Limit {
	rdl := &RedisLimit{
		rc: cli,
	}
	for _, ll := range hdl {
		ll(rdl)
	}
	l := &Limit{Key: name, L: rdl}
	limitManager.Store(name, l)
	return l
}

func GetLimit(name string) *Limit {
	l, ok := limitManager.Load(name)
	if !ok {
		return nil
	}
	return l.(*Limit)
}

func (l *Limit) DoLimit(current int64, key, value string) (isLimit bool, err error) {
	limit := GetLimit(key)
	if limit == nil {
		err = fmt.Errorf("invalid limiter")
		return
	}
	start := limit.L.CalcStart(current)
	defer func() {
		if current%10 < 2 {
			_, err = limit.L.Del(start)
			if err != nil {
				fmt.Printf("del error:%v", err)
			}
		}
	}()
	_, err = limit.L.Add(value, current)
	if err != nil {
		err = fmt.Errorf("invalid limiter:%v", err)
		return
	}
	isLimited, err := limit.L.IsLimit(start, current)
	if err != nil {
		err = fmt.Errorf("invalid limiter:%v", err)
		return
	}
	if isLimited {
		return true, nil
	}
	return
}

func (rl *RedisLimit) Add(value string, timeStamp int64) (int64, error) {
	return rl.rc.ZAdd(rl.Key, redis.Z{Score: float64(timeStamp), Member: value}).Result()
}

func (rl *RedisLimit) IsLimit(start, end int64) (bool, error) {
	startStr := strconv.FormatInt(start, 10)
	endStr := strconv.FormatInt(end, 10)
	count, err := rl.rc.ZCount(rl.Key, startStr, endStr).Result()
	if err != nil {
		return false, err
	}
	if count >= rl.Count {
		return true, nil
	}
	return false, nil
}

func (rl *RedisLimit) Del(start int64) (int64, error) {
	startStr := strconv.FormatInt(start, 10)
	return rl.rc.ZRemRangeByScore(rl.Key, "-inf", startStr).Result()
}

func (rl *RedisLimit) CalcStart(current int64) int64 {
	return current - rl.Duration
}
