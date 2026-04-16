package database

import (
	"homedy/config"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Password: config.REDIS_PASS,
		DB:       0,
	})
}
