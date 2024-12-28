package main

import (
	"log"
	"os"

	"byteload-agent/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	srv := server.New(port)

	log.Printf("Server starting on port %s", port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
