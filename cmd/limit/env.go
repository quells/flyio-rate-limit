package main

import (
	"os"
	"strconv"
)

func envInt(key string, fallback int) int {
	x, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}

	return x
}
