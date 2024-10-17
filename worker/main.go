package main

import (
	"log"
	"os"

	// rd => redis driver
	rd "github.com/redis/go-redis/v9"
)

var redis *rd.Client

func init() {
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

}
