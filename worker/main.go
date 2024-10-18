package main

import (
	"context"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	rd "github.com/redis/go-redis/v9" // rd => redis driver
)

var is_dev bool
var redis = RedisDefault()
var job_id_idx = atomic.Uint32 {}

func init() {
	is_dev = len(os.Args) == 2 && os.Args[1] == "-dev"

	if is_dev {
		if err := godotenv.Load("../.env.public", "../.env"); err != nil {
			log.Fatalf("Failed to load ../.env.public or ../.env environment file: %s\n", err)
		}
	}

	redis.set(rd.NewClient(&rd.Options{
		Addr:     get_env("REDIS_ADDR"),
		Password: get_env("REDIS_PASSWORD"),
	}))
}

func main() {
	ctx := context.Background()

	go func() {
		pubsub := redis.subscribe(ctx, get_env("REDIS_CH_WORK_PROCESS"))
		defer pubsub.Close()

		for message := range pubsub.Channel() {
			go process_work(message.Payload)
		}
	}()

	ServeForever()
}

func process_work(s string) {
    // GET config data via REDIS_WORK_PROCESS
    // DEL the key
    // process
    // SET REDIS_WORK_RESULT process_data
	log.Println(s)
	panic("unimplemented")
}

// Increments and returns the current job index
func new_job_idx() uint32 {
    var new_val uint32

    for {
        val := job_id_idx.Load()
        new_val = (val + 1) % 0xFF // Between 1 and 255 seems alright since it shouldn't have that much jobs overlapping
        if new_val == 0 {
            new_val = 1
        }

        if job_id_idx.CompareAndSwap(val, new_val) {
            break
        }
    }

    return new_val
}

func get_env(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Panicf("Error: env var not found for key: '%s'\n", key)
	}

	return value
}
