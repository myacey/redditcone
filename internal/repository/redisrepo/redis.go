package redisrepo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func ConfigureRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       0,

		// pool settings
		PoolSize:     10,
		MinIdleConns: 5,
		MaxConnAge:   30 * time.Minute,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %v", err)
	}

	return rdb, nil
}
