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

	if envAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		serverAddr = envAddr
	}
	if envBase, ok := os.LookupEnv("BASE_URL"); ok {
		baseURL = envBase
	}
	if envLog, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel = envLog
	}
	if envPath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		fileStoragePath = envPath
	}

	return &Config{
		ServerAddr:      serverAddr,
		BaseURL:         baseURL,
		LogLevel:        logLevel,
		FileStoragePath: fileStoragePath,
	}
}
