package main

import (
	"log"
	"shortUrl/shorten_url/cmd/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
