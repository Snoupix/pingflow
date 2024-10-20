package main_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	worker "worker"
	. "worker/utils"

	"github.com/joho/godotenv"
	rd "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const (
	WORKER_ADDR           = "WORKER_ADDR"
	WORKER_PORT           = "WORKER_PORT"
	REDIS_PORT            = "REDIS_PORT"
	REDIS_ADDR            = "REDIS_ADDR"
	REDIS_PASSWORD        = "REDIS_PASSWORD"
	REDIS_WORK_PREFIX     = "REDIS_WORK_PREFIX"
	REDIS_WORK_PROCESS    = "REDIS_WORK_PROCESS"
	REDIS_WORK_RESULT     = "REDIS_WORK_RESULT"
	REDIS_WORK_STATUS     = "REDIS_WORK_STATUS"
	REDIS_CH_WORK_PROCESS = "REDIS_CH_WORK_PROCESS"
	REDIS_CH_WORK_RESULT  = "REDIS_CH_WORK_RESULT"
)

func initEnv(t *testing.T) {
	err := godotenv.Load("../.env.public", "../.env")
	assert.Nil(t, err, "Failed to load ../.env.public or ../.env environment file: %v", err)

	val, err := TryGetEnv(WORKER_ADDR)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(WORKER_PORT)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_PORT)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_ADDR)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_WORK_PREFIX)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_WORK_PROCESS)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_WORK_RESULT)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_WORK_STATUS)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_CH_WORK_PROCESS)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")

	val, err = TryGetEnv(REDIS_CH_WORK_RESULT)
	assert.Nil(t, err, err)
	assert.NotEqual(t, strings.Trim(val, " "), "")
}

func fetchJobIndex(t *testing.T, httpclient *http.Client) string {
	url := fmt.Sprintf("http://%s:%s/job-index", GetEnv(WORKER_ADDR), GetEnv(WORKER_PORT))
	resp, err := httpclient.Get(url)
	assert.Nil(t, err, fmt.Sprintf("Cannot access the worker REST API /job-idx endpoint %s", url))
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, 200, "/job-index response status should be 200")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, fmt.Sprintf("Cannot parse the /job-idx response body %s", err))
	assert.NotEmpty(t, body, "/job-index shouldn't be empty")

	return string(body)
}

func TestWorker(t *testing.T) {
	initEnv(t)

	ctx := context.Background()
	httpclient := http.Client{}
	redis := worker.RedisDefault()

	job_idx := fetchJobIndex(t, &httpclient)

	redis.Set(rd.NewClient(&rd.Options{
		Addr:     fmt.Sprintf("%s:%s", GetEnv(REDIS_ADDR), GetEnv(REDIS_PORT)),
		Password: GetEnv(REDIS_PASSWORD),
	}))

	client := redis.Get()
	defer redis.Unlock()

	config_key := fmt.Sprintf("%s:%s:%s", GetEnv(REDIS_WORK_PREFIX), job_idx, GetEnv(REDIS_WORK_PROCESS))
	err := client.HSet(ctx, config_key, "endpoint", "/api", "parameters", "classes").Err()
	assert.Nil(t, err, err)

	sub := client.Subscribe(ctx, GetEnv(REDIS_CH_WORK_RESULT))
	defer sub.Close()

	sub_mess, err := sub.Receive(ctx)
	assert.Nil(t, err, "Failed to receive message from subscriber channel (Subscription step)")

	err = client.Publish(ctx, GetEnv(REDIS_CH_WORK_PROCESS), job_idx).Err()
	assert.Nil(t, err, err)

	var result_key string

	// This might fail since it's not concurrently subscribed to the channel before triggering the work
	switch sub_mess.(type) {
	case *rd.Subscription:
		sub_mess, err = sub.Receive(ctx)
		assert.Nil(t, err, "Failed to receive message from subscriber channel (Message step)")

		switch sub_mess.(type) {
		case *rd.Message:
			result_key = sub_mess.(*rd.Message).Payload
		default:
			t.Log("Subscription channel: Failed to receive Message response")
			t.Fail()
			return
		}

	default:
		t.Log("Subscription channel: Failed to receive Subscription response")
		t.Fail()
		return
	}

	result, err := client.Get(ctx, result_key).Result()
	assert.Nil(t, err, err)

	assert.NotEqual(t, len(result), 0, "Work result string shouldn't be empty")

	// For visual debug
	// fmt.Println(result)
}
