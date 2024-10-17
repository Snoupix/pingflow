package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	rd "github.com/redis/go-redis/v9" // rd => redis driver
)

var is_dev bool
var redis *rd.Client

func init() {
    is_dev = len(os.Args) == 2 && os.Args[1] == "-dev"

    if is_dev {
        if err := godotenv.Load("../.env.public"); err != nil {
            log.Fatalf("Failed to load ../.env.public env file: %s\n", err)
        }
    }

    redis_addr_key := "REDIS_ADDR"
    redis_addr, ok := os.LookupEnv(redis_addr_key)
    if !ok {
        log.Panicf("Error: redis addr not found in env. Key: '%s'\n", redis_addr_key)
    }

    redis = rd.NewClient(&rd.Options{
        Addr: redis_addr,
    })
}

func main() {
    ServeForever(is_dev)
}
