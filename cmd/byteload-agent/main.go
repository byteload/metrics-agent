package main

import (
	"log"
	"os"

	"byteload-agent/internal/server"

	"github.com/spf13/viper"
)

// Version is set during build
var version = "dev"

const (
	defaultPort = "9001"
)

func main() {
	// Set version in server package
	server.Version = version

	configFile := os.Getenv("BYTELOAD_CONFIG_FILE")
	if configFile == "" {
		configFile = "/etc/byteload/byteload.yaml"
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %s, using defaults", err)
	}

	port := viper.GetString("server.port")
	if port == "" {
		port = defaultPort
	}

	srv := server.New(server.Config{
		Port: port,
		Auth: server.AuthConfig{
			Enabled:  viper.GetBool("security.basic_auth.enabled"),
			Username: viper.GetString("security.basic_auth.username"),
			Password: viper.GetString("security.basic_auth.password"),
		},
		Services: server.ServicesConfig{
			Default: viper.GetStringSlice("services.default"),
		},
	})

	log.Printf("Starting byteload agent version %s on port %s", version, port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
