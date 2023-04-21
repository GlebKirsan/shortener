package config

import (
	"flag"
)

var (
	ServerAddress  = flag.String("a", "localhost:8080", "address and port of the server")
	ResponsePrefix = "http://localhost:8080"
)

func init() {
	flag.Func("b", "Response prefix", func(flagValue string) error {
		if flagValue[len(flagValue)-1] == '/' {
			ResponsePrefix = flagValue[:len(flagValue)-1]
		} else {
			ResponsePrefix = flagValue
		}
		return nil
	})
}
