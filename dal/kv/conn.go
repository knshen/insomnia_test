package kv

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

var (
	redisCli *redis.Client
)

func InitRedisClient() {
	miniRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	redisCli = redis.NewClient(&redis.Options{
		Addr:     miniRedis.Addr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

func GetRedisClient() *redis.Client {
	return redisCli
}
