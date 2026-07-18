package config

import "flag"

type Config struct {
	FlagRunAddr  string
	FlagHTTPAddr string
}

func Load() *Config {
	var runAddr, httpAddr string
	flag.StringVar(&runAddr, "a", "0.0.0.0:8080", "address and port to run server")
	flag.StringVar(&httpAddr, "b", "http://0.0.0.0:8080", "address server")
	flag.Parse()

	return &Config{
		FlagRunAddr:  runAddr,
		FlagHTTPAddr: httpAddr,
	}
}
