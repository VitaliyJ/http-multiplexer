package app

import "os"

type Config struct {
	HTTPServerPort string
	LogLevel       string
}

func NewConfig() *Config {
	return &Config{
		HTTPServerPort: getOrDefault("HTTP_SERVER_PORT", "8080"),
		LogLevel:       getOrDefault("LOG_LEVEL", "info"),
	}
}

func getOrDefault(env string, d string) string {
	v := os.Getenv(env)
	if v == "" {
		return d
	}

	return v
}
