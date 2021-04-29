package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func getConnection() {
	if rdb == nil {
		fmt.Println("Connecting to redis")
		rdb = redis.NewClient(&redis.Options{
			Addr:     "192.168.86.250:32768",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
}

func Set(key string, value string, expiration time.Duration) error {
	getConnection()

	// store the value in redis
	err := rdb.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) string {
	getConnection()

	value, err := rdb.Get(key).Result()

	if err != nil {
		return ""
	}

	return value
}
