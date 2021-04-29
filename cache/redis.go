package cache

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var rdb *redis.Client
var redisHost string
var redisPort int

func getConnection() {
	if rdb == nil {
		var err error
		redisHost = os.Getenv("redisHost")
		redisPort, err = strconv.Atoi(os.Getenv("redisPort"))

		if err != nil {
			fmt.Println("failed to get redis vals")
		}

		fmt.Println("Connecting to redis")
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + strconv.Itoa(redisPort),
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
