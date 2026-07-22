package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr  string
	BaseURL     string
	LogLevel    string
	DatabaseDSN string
}

func Load() *Config {
	var serverAddr, baseURL, logLevel, databaseDSN string
	flag.StringVar(&serverAddr, "a", "0.0.0.0:8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://0.0.0.0:8080", "address server")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.StringVar(&databaseDSN, "d", "postgres://user:password@localhost:5432/database", "database DSN")
	flag.Parse()

	if envAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		serverAddr = envAddr
	}
	if envBase, ok := os.LookupEnv("BASE_URL"); ok {
		baseURL = envBase
	}
	if envLog, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel = envLog
	}
	if envDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		databaseDSN = envDSN
	}

	return &Config{
		ServerAddr:  serverAddr,
		BaseURL:     baseURL,
		LogLevel:    logLevel,
		DatabaseDSN: databaseDSN,
	}
}
