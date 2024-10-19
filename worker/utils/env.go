package utils

import (
	"log"
	"os"
)

var env_cache = map[string]string{}

func GetEnv(key string) string {
    if value, ok := env_cache[key]; ok {
        return value
    }

    value, ok := os.LookupEnv(key)
	if !ok {
		log.Panicf("Error: env var not found for key: '%s'\n", key)
	}

    env_cache[key] = value

	return value
}
