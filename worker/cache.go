package main

import "context"

const CACHE_KEY = "cache" // Only used by the worker

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

func (c *Cache) Store(ctx context.Context, redis *Redis, key string, value string) error {
	client := redis.Get()
	defer redis.Unlock()

	if err := client.HSet(ctx, key, value).Err(); err != nil {
		return err
	}

	(*c)[key] = value

	return nil
}
