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

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var param RequestBody
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	shortenedUrl, err := h.service.Create(param.Original)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, shortenedUrl)
}
