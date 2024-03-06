package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortUrl/shorten_url/internal/service"
)

type URLHandler struct {
	service service.URLService
}

type RequestBody struct {
	Original string `json:"original"`
}

func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var param RequestBody
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortenedUrl, err := h.service.Create(param.Original)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shortened_url": shortenedUrl.Shortened})
}

func (h *URLHandler) GetURL(c *gin.Context) {
	shortened := c.Param("shortened")
	if shortened == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not provided"})
		return
	}

	url, err := h.service.Get(shortened)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}
	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"original_url": url.Original})
}
