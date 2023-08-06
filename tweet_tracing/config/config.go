package config

import "os"

type Config struct {
	Address       string
	JaegerAddress string
}

func GetConfig() Config {
	return Config{
		Address:       os.Getenv("TWEET_TRACING_SERVICE_ADDRESS"),
		JaegerAddress: os.Getenv("JAEGER_ADDRESS"),
	}
}
