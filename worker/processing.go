package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"worker/utils"
)

type WorkConfig struct {
	endpoint   string
	parameters string // Note that it can be a space separeted list of parameters
}

func ProcessWork(ctx context.Context, httpclient *http.Client, work_id string) {
	client := redis.Get()
	defer redis.Unlock() // FIXME: Is this Mutex guard useful ?

	config_key := fmt.Sprintf("%s:%s:%s", utils.GetEnv("REDIS_WORK_PREFIX"), work_id, utils.GetEnv("REDIS_WORK_PROCESS"))

	cmd := client.HGetAll(ctx, config_key) // Might wanna handle rd.Nil result ?
	res, err := cmd.Result()
	if err != nil {
		log.Fatalf("Failed to get config for work ID: %s. Config key: %s", work_id, config_key)
	}

	endpoint, ok := res["endpoint"]
	if !ok || len(endpoint) == 0 {
		log.Fatalf("Malformed config data, missing or empty 'endpoint' field")
	}
	parameters, ok := res["parameters"]
	if !ok {
		log.Fatalf("Malformed config data, missing 'parameters' field")
	}

	config := WorkConfig{
		endpoint,
		parameters,
	}

	client.Del(ctx, config_key)

	output, ok := cache[config.endpoint]
	if !ok {
		output, ok = FetchEndpoint(ctx, httpclient, config)

		if ok {
			// Caching a request result with "endpointparameters" concatenated as a key so
			// the same request has the same result
			cache.Store(ctx, client, config.endpoint+config.parameters, output)
		}
	}

	result_key := fmt.Sprintf("%s:%s:%s", utils.GetEnv("REDIS_WORK_PREFIX"), work_id, utils.GetEnv("REDIS_WORK_RESULT"))
	client.Set(ctx, result_key, output, 0) // TODO: Set TTL ?

	client.Publish(ctx, utils.GetEnv("REDIS_CH_WORK_RESULT"), result_key)
}
