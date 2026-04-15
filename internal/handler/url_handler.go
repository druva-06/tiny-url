package handler

import (
	"net/http"

	"github.com/druva-06/tiny-url/internal/dto/request"
	"github.com/druva-06/tiny-url/internal/service"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{service: s}
}

func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req request.CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil { //
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}
	shortCode, err := h.service.CreateShortURL(req.LongUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"short_code": shortCode})
}

func (h *URLHandler) GetLongURL(c *gin.Context) {
	code := c.Param("code")
	longUrl, err := h.service.GetLongURL(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"long_url": longUrl})
}
