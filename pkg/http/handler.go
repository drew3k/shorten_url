package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortUrl/shorten_url/internal/service"
	"strconv"
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
	var param *RequestBody
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

func (h *URLHandler) RedirectURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not provided"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	url, err := h.service.Get(idInt)
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
