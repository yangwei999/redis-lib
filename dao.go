package redisdb

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var errKeyNotExists = errors.New("key doesn't exist")

func (cli *client) RPush(key string, val interface{}) error {
	return cli.withContext(func(ctx context.Context) error {
		return cli.redisCli.RPush(ctx, key, val).Err()
	})
}

func (cli *client) LPop(key string, val interface{}) error {
	return cli.withContext(func(ctx context.Context) error {
		err := cli.redisCli.LPop(ctx, key).Scan(val)
		if err == redis.Nil {
			return errKeyNotExists
		}

		return err
	})
}

func (cli *client) Set(key string, val interface{}) error {
	return cli.withContext(func(ctx context.Context) error {
		return cli.redisCli.Set(ctx, key, val, 0).Err()
	})
}

func (cli *client) SetWithExpiry(key string, val interface{}, expiry time.Duration) error {
	return cli.withContext(func(ctx context.Context) error {
		return cli.redisCli.Set(ctx, key, val, expiry).Err()
	})
}

func (cli *client) Get(key string, data interface{}) error {
	return cli.withContext(func(ctx context.Context) error {
		err := cli.redisCli.Get(ctx, key).Scan(data)
		if err == redis.Nil {
			return errKeyNotExists
		}

		return err
	})
}

func (cli *client) Expire(key string, expire time.Duration) error {
	return cli.withContext(func(ctx context.Context) error {
		return cli.redisCli.Expire(ctx, key, expire).Err()
	})
}

func (cli *client) SetKey(key string, expiry time.Duration) error {
	return cli.withContext(func(ctx context.Context) error {
		return cli.redisCli.Set(ctx, key, 0, expiry).Err()
	})
}

func (cli *client) HasKey(key string) (bool, error) {
	exists := false
	err := cli.withContext(func(ctx context.Context) error {
		n, err := cli.redisCli.Exists(ctx, key).Result()
		if err != nil {
			return err
		}

		exists = n > 0

		return nil
	})

	return exists, err
}

func (cli *client) IsKeyNotExists(err error) bool {
	return errors.Is(err, errKeyNotExists)
}
