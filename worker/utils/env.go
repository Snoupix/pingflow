package utils

import (
	"log"
	"os"
)

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Panicf("Error: env var not found for key: '%s'\n", key)
	}

	return value
}
