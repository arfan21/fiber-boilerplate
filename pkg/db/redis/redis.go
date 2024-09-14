package dbredis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arfan21/fiber-boilerplate/config"
	"github.com/arfan21/fiber-boilerplate/pkg/logger"
	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Get().Redis.URL, config.Get().Redis.Port),
		Username: config.Get().Redis.Username,
		Password: config.Get().Redis.Password,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		logger.Log(context.Background()).Error().Err(err).Msg("failed to ping redis")
		return nil, err
	}

	logger.Log(context.Background()).Info().Msg("dbredis: connection established")

	return client, nil
}

func Get[T any](ctx context.Context, conn *redis.Client, key string) (res T, err error) {
	val, err := conn.Get(ctx, key).Bytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(val, &res)

	return res, err
}

func Set[T any](ctx context.Context, conn *redis.Client, key string, value T, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return conn.Set(ctx, key, string(val), expiration).Err()
}
