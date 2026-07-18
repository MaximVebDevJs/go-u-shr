package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr      string
	BaseURL         string
	LogLevel        string
	FileStoragePath string
}

func Load() *Config {
	var serverAddr, baseURL, logLevel, fileStoragePath string
	flag.StringVar(&serverAddr, "a", "0.0.0.0:8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://0.0.0.0:8080", "address server")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.StringVar(&fileStoragePath, "f", "storage.json", "file storage path")
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
	if envPath := os.Getenv("FILE_STORAGE_PATH"); envPath != "" {
		fileStoragePath = envPath
	}

	return &Config{
		ServerAddr:      serverAddr,
		BaseURL:         baseURL,
		LogLevel:        logLevel,
		FileStoragePath: fileStoragePath,
	}
}
