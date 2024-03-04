package server

import (
	"github.com/gin-gonic/gin"
	"shortUrl/shorten_url/internal/repository"
	"shortUrl/shorten_url/internal/service"
	"shortUrl/shorten_url/pkg/http"
)

func Run() error {
	r := gin.Default()

	urlRepo := repository.NewInMemoryURLRepository()
	urlService := service.NewUrlRepository(urlRepo)
	urlHandler := http.NewURLHandler(urlService)

	r.POST("/shorten", urlHandler.ShortenURL)
	r.GET("/shorten/:shortened", urlHandler.GetURL)

	return r.Run(":8080")
}
