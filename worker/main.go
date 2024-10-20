package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	rd "github.com/redis/go-redis/v9" // rd => redis driver

	. "worker/utils"
)

var ctx = context.Background()
var is_dev bool
var redis = RedisDefault()
var cache Cache

func main() {
	httpclient := http.Client{}

	is_dev = len(os.Args) == 2 && os.Args[1] == "-dev"

	if is_dev {
		if err := godotenv.Load("../.env.public", "../.env"); err != nil {
			log.Fatalf("Failed to load ../.env.public or ../.env environment file: %s\n", err)
		}
	}

	redis.Set(rd.NewClient(&rd.Options{
		Addr:     fmt.Sprintf("%s:%s", GetEnv("REDIS_ADDR"), GetEnv("REDIS_PORT")),
		Password: GetEnv("REDIS_PASSWORD"),
	}))

	if err := cache.Load(ctx, &redis); err != nil {
		log.Fatalf("Failed to load cache from Redis: %s\n", err)
	}

	go func() {
		pubsub := redis.Subscribe(ctx, GetEnv("REDIS_CH_WORK_PROCESS"))
		defer pubsub.Close()

		for message := range pubsub.Channel() {
			go ProcessWork(ctx, &httpclient, message.Payload)
		}
	}()

	ServeForever()
}
