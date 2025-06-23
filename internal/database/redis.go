package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"insider_task/internal/configs"
)

func ConnectRedis(ctx context.Context, config *configs.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
	})
	_, err := client.Ping(ctx).Result()
	return client, err
}
