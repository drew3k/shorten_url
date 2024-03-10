package main

import (
	"log"
	"shortUrl/shorten_url/cmd/server"
)

func main() {
	s := server.NewServer()
	if err := s.Run(); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
