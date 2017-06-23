package app

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func envInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Panic(err)
		}
		return int(v)
	}
	return defaultValue
}

func envDuration(key, defaultValue string) Interval {
	value, err := time.ParseDuration(env(key, defaultValue))
	if err != nil {
		log.Panic(err)
	}
	return Interval(value)
}

func envBool(key, defaultValue string) bool {
	return "true" == env(key, defaultValue)
}

func envStringList(key, defaultValue string) StringList {
	if value := env(key, defaultValue); value != "" {
		return StringList(strings.Split(value, ";"))
	}
	return StringList([]string{})
}

func envUint32(key string, defaultValue int) uint32 {
	return uint32(envInt(key, defaultValue))
}
