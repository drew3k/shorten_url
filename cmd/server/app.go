package server

import (
	"github.com/gin-gonic/gin"
	"shortUrl/shorten_url/internal/domain"
	"shortUrl/shorten_url/internal/repository"
	"shortUrl/shorten_url/internal/service"
	"shortUrl/shorten_url/internal/telegram"
	"shortUrl/shorten_url/pkg/http"
)

type Server struct {
	router       *gin.Engine
	shortenedURL *domain.URL
	bot          telegram.BotService
}

func NewServer(bot telegram.BotService) *Server {
	return &Server{
		router: gin.Default(),
		bot:    bot,
	}
}

func (s *Server) SetupRoutes() {
	urlRepo := repository.NewInMemoryURLRepository()
	urlService := service.NewUrlService(urlRepo)
	urlHandler := http.NewURLHandler(urlService)

	s.router.POST("/shorten", urlHandler.ShortenURL)
}

func (s *Server) Run() error {
	if err := s.bot.Initialize(); err != nil {
		return err
	}
	s.SetupRoutes()
	go s.bot.StartTelegramBot()
	return s.router.Run(":8080")
}
