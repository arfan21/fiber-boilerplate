package userrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arfan21/fiber-boilerplate/internal/entity"
	"github.com/arfan21/fiber-boilerplate/pkg/constant"
	dbredis "github.com/arfan21/fiber-boilerplate/pkg/db/redis"
	"github.com/redis/go-redis/v9"
)

type RepositoryRedis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *RepositoryRedis {
	return &RepositoryRedis{client: client}
}

func (r RepositoryRedis) SetRefreshToken(ctx context.Context, token string, expireIn time.Duration, payload entity.UserRefreshToken) (err error) {
	fmt.Println("expireIn", expireIn)
	key := constant.RedisRefreshTokenKeyPrefix + token
	err = dbredis.Set(ctx, r.client, key, payload, expireIn)
	if err != nil {
		err = fmt.Errorf("user.repository_redis.SetRefreshToken: failed to set refresh token: %w", err)
		return
	}

	return
}

func (r RepositoryRedis) IsRefreshTokenExist(ctx context.Context, token string) (payload entity.UserRefreshToken, err error) {
	key := constant.RedisRefreshTokenKeyPrefix + token
	payload, err = dbredis.Get[entity.UserRefreshToken](ctx, r.client, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = constant.ErrUnauthorizedAccess
		}
		err = fmt.Errorf("user.repository_redis.IsRefreshTokenExist: failed to get refresh token: %w", err)
		return
	}
	return
}

func (r RepositoryRedis) DeleteRefreshToken(ctx context.Context, token string) (err error) {
	key := constant.RedisRefreshTokenKeyPrefix + token
	err = r.client.Del(ctx, key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		err = fmt.Errorf("user.repository_redis.DeleteRefreshToken: failed to delete refresh token: %w", err)
		return
	}

	return
}
