package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	rd "github.com/redis/go-redis/v9" // rd => redis driver
)

var is_dev bool
var redis = RedisDefault()

func init() {
	is_dev = len(os.Args) == 2 && os.Args[1] == "-dev"

	if is_dev {
		if err := godotenv.Load("../.env.public"); err != nil {
			log.Fatalf("Failed to load ../.env.public env file: %s\n", err)
		}
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatalf("Failed to load ../.env env file: %s\n", err)
		}
	}

	redis.set(rd.NewClient(&rd.Options{
		Addr:     getEnv("REDIS_ADDR"),
		Password: getEnv("REDIS_PASSWORD"),
	}))
}

func main() {
	ctx := context.Background()

	go func() {
		channel := redis.subscribe(ctx, "work:process")

		for message := range channel {
			go process_work(*&message.Payload)
		}
	}()

	ServeForever(is_dev)
}

func process_work(s string) {
	panic("unimplemented")
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Panicf("Error: env var not found for key: '%s'\n", key)
	}

	return value
}
