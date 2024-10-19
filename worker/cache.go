package main

import (
	"context"

	rd "github.com/redis/go-redis/v9" // rd => redis driver
)

const CACHE_KEY = "work:cache" // Only used by the worker

type Cache map[string]string // maps Endpoint => JSON parsed Data

func (c *Cache) Load(ctx context.Context, redis *Redis) error {
	client := redis.Get()
	defer redis.Unlock()

	res := client.HGetAll(ctx, CACHE_KEY)
	err := res.Err()
	if err != nil {
		return err
	}

	*c, err = res.Result()

	return err
}

func (c *Cache) Store(ctx context.Context, client *rd.Client, key string, value string) error {
	if err := client.HSet(ctx, CACHE_KEY, key, value).Err(); err != nil {
		return err
	}

	(*c)[key] = value

	return nil
}
