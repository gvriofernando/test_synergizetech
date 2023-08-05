package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Address  string
	Password string
	Database int
}

func NewClient(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		PoolSize:        1000,
		Addr:            cfg.Address,
		Password:        cfg.Password,
		DB:              cfg.Database,
		MaxRetries:      6,
		MaxRetryBackoff: 5 * time.Second,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
