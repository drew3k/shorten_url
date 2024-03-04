package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortUrl/shorten_url/internal/service"
)

type URLHandler struct {
	service service.URLService
}

var requestBody struct {
	Original string `json:"original"`
}

func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortenedUrl, err := h.service.Create(requestBody.Original)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shortened_url": shortenedUrl.Shortened})
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	var id int

	url, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}
	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(http.StatusPermanentRedirect, url.Original)
}
