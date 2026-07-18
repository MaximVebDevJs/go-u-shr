package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr string
	BaseURL    string
	LogLevel   string
}

func Load() *Config {
	var serverAddr, baseURL, logLevel string
	flag.StringVar(&serverAddr, "a", "0.0.0.0:8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://0.0.0.0:8080", "address server")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		serverAddr = envAddr
	}
	if envBase := os.Getenv("BASE_URL"); envBase != "" {
		baseURL = envBase
	}
	if envLog := os.Getenv("LOG_LEVEL"); envLog != "" {
		logLevel = envLog
	}

	return &Config{
		ServerAddr: serverAddr,
		BaseURL:    baseURL,
		LogLevel:   logLevel,
	}
}
