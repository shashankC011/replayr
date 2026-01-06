package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	CAPTURE_DIR_ADDR string
	REPLAYS_DIR_ADDR string
	PORT             string
	SERVER_ADDR      string
}

func Load() *Config {
	cfg := &Config{
		CAPTURE_DIR_ADDR: getEnv("CAPTURE_DIR_ADDR", "./"),
		REPLAYS_DIR_ADDR: getEnv("REPLAYS_DIR_ADDR", "./"),
		PORT:             getEnv("PORT", "9000"),
		SERVER_ADDR:      getEnv("SERVER_ADDR", "http://localhost:8080"),
	}
	return cfg
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("Failed to load env variable: ", key)
	}
	return val
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	fmt.Printf("Couldnt find value for env variable: %s therefore using fallback: %s", key, fallback)
	return fallback
}
