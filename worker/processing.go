package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"worker/utils"
)

type WorkConfig struct {
	endpoint   string
	parameters string // Note that it can be a space separeted list of parameters
}

type ColorOut struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
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
        // If ok is false, it just means that there is a config parsing error, else it will panic

		if ok {
			// Caching a request result with "endpointparameters" concatenated as a key so
			// the same request has the same result
			cache.Store(ctx, client, config.endpoint+config.parameters, output)
		}
	}

	result_key := fmt.Sprintf("%s:%s:%s", utils.GetEnv("REDIS_WORK_PREFIX"), work_id, utils.GetEnv("REDIS_WORK_RESULT"))
	client.Set(ctx, result_key, output, time.Second*30)

	client.Publish(ctx, utils.GetEnv("REDIS_CH_WORK_RESULT"), result_key)
}

func ProcessColor(ctx context.Context, httpclient *http.Client) {
	client := redis.Get()
	defer redis.Unlock()

	// Values between 100 and 255 for more brightness
	r := uint8(rand.Intn(156) + 100)
	g := uint8(rand.Intn(156) + 100)
	b := uint8(rand.Intn(156) + 100)

	data, err := json.Marshal(ColorOut{r, g, b})
	if err != nil {
		log.Fatalf("Unreachable JSON color parsing")
	}

	client.Publish(ctx, utils.GetEnv("REDIS_CH_COLOR_RESULT"), string(data))
}
