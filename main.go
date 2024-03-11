package main

import (
	"log"
	"shortUrl/shorten_url/cmd/server"
	"shortUrl/shorten_url/internal/telegram"
)

func main() {
	bot := telegram.NewBotAPI(nil)
	err := bot.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize bot:", err)
	}

	s := server.NewServer(bot)
	if err := s.Run(); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
