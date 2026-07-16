package config

import "flag"

var FlagRunAddr string
var FlagHTTPAddr string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "0.0.0.0:8080", "address and port to run server")
	flag.StringVar(&FlagHTTPAddr, "b", "http://0.0.0.0:8080", "address server")
	flag.Parse()
}
