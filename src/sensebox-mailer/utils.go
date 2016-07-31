package main

import (
	"fmt"
	"os"
)

func getBytesFromEnvOrFail(key string) []byte {
	envBytes := []byte(os.Getenv(ENV_PREFIX + key))
	if len(envBytes) == 0 {
		fmt.Println("Please add", ENV_PREFIX+key, "to your environment")
		os.Exit(1)
		return nil
	} else {
		return envBytes
	}
}
