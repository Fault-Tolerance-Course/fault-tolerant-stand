package redis

import (
	"context"

	"ad-service/config"
	"ad-service/internal/pkg/closer"

	"github.com/redis/go-redis/v9"
)

func Client(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        config.Instance().Redis.Addr,
		DB:          config.Instance().Redis.DB,
		MaxRetries:  config.Instance().Redis.MaxRetries,
		DialTimeout: config.Instance().Redis.DialTimeout,
		PoolTimeout: config.Instance().Redis.Timeout,
	})

	closer.Add(client.Close)

	cmd := client.Ping(ctx)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return client, nil
}
