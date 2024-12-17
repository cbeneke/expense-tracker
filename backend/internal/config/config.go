package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	return Config{
		Port: getEnvWithDefault("PORT", "8080"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
