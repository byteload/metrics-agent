package main

import (
	"log"
	"os"

	"byteload_agent/internal/server"

	"github.com/spf13/viper"
)

const (
	defaultPort = "9001"
)

func main() {
	configFile := os.Getenv("BYTELOAD_CONFIG_FILE")
	if configFile == "" {
		configFile = "/etc/byteload/byteload.yaml"
	}

	viper.SetConfigFile(configFile)
	port := defaultPort
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %s, using default port %s", err, defaultPort)
	} else {
		if p := viper.GetString("server.port"); p != "" {
			port = p
		}
	}

	srv := server.New(port)

	log.Printf("Server starting on port %s", port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
