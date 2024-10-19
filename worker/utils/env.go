package utils

import (
	"fmt"
	"os"
)

var env_cache = map[string]string{}

func TryGetEnv(key string) (string, error) {
    if value, ok := env_cache[key]; ok {
        return value, nil
    }

    value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("Error: env var not found for key: '%s'\n", key)
	}

    env_cache[key] = value

	return value, nil
}

// This function panics if key is not found.
// Use TryGetEnv if you want to handle the result
func GetEnv(key string) string {
    val, err := TryGetEnv(key)
    if err != nil {
        panic(err)
    }

    return val
}
