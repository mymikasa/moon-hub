package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"moon/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrKeyNotExist = redis.Nil

// go:generate mockgen -source=user.go -package=cachemocks -destination=./mocks/user.mock.go UserCache
type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, uid int64, u domain.User) error
	Del(ctx context.Context, uid int64) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)

	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	if err != nil {
		return domain.User{}, err
	}
	return u, err
}

func (c *RedisUserCache) Set(ctx context.Context, uid int64, u domain.User) error {
	key := c.key(uid)

	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisUserCache) Del(ctx context.Context, uid int64) error {
	key := c.key(uid)
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
