package limit

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/transfercai/redis-extension/client"
)

func TestDoLimit(t *testing.T) {
	name := "test"
	l := NewRedisLimit(name, client.DoInjectRedis(redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       0,
	})), SetCount(100), SetDuration(1000))
	for i := 0; i < 1000; i++ {
		cur := time.Now().UnixNano() / 1e6
		isLimit, err := l.DoLimit(cur, name, fmt.Sprint(cur))
		if err != nil {
			t.Fatal(err)
		}
		if isLimit {
			t.Logf("num is:%d", i)
		}
		time.Sleep(time.Millisecond * 2)
	}

}
