package redis_extension

import (
	"context"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     ":6379", // use default Addr
		Password: "",      // no password set
		DB:       0,       // use default DB
	})
}

func TestExtensionIncrAndExpire(t *testing.T) {
	c := DoInjectRedis(getRedisClient())
	_, err := c.IncrAndExpire(context.TODO(), "333", 1000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtensionSetNxAndExpire(t *testing.T) {
	c := DoInjectRedis(getRedisClient())
	_, err := c.SetNxAndExpire(context.TODO(), "3333", 1001, 1000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtensionGetAndDel(t *testing.T) {
	c := DoInjectRedis(getRedisClient())
	_, err := c.GetAndDel(context.TODO(), "3333", 1001)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtensionDecrAndExpire(t *testing.T) {
	c := DoInjectRedis(getRedisClient())
	_, err := c.DecrAndExpire(context.TODO(), "3333", 1001)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtensionHSetAndExpire(t *testing.T) {
	c := DoInjectRedis(getRedisClient())
	_, err := c.HSetAndExpire(context.TODO(), "1111", "child", 1001, 1000)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkExtensionSetNxAndExpire(b *testing.B) {
	c := DoInjectRedis(getRedisClient())
	for i := 0; i < b.N; i++ {
		//_, err := c.SetNX("3333", 1001, time.Second*1000).Result()
		//_, err := c.SetNxAndExpire(context.TODO(), "3333", 1001, 1000)
		_, err := c.Set("222", 1001, time.Second*10).Result()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestExtensionMSetAndExpire(t *testing.T) {
	keys := []string{"12", "123", "1234", "r3", "334"}
	c := DoInjectRedis(getRedisClient())
	_, err := c.MSetAndExpire(context.TODO(), keys, []interface{}{1, "3", 4, 5, 6}, 1002)
	if err != nil {
		t.Fatal(err)
	}
}
