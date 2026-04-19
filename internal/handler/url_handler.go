package handler

import (
	"context"
	"log"
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
	longUrl, err := h.service.GetLongURL(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"long_url": longUrl})
}

func (h *URLHandler) UpdateLongUrl(c *gin.Context) {
	code := c.Param("code")
	var req request.CreateShortURLRequest
	log.Printf("[UpdateLongUrl] START code=%s", code)
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateLongUrl] ERROR code=%s error=%s", code, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[UpdateLongUrl] REQUEST code=%s request=%+v", code, req)
	res, err := h.service.UpdateLongUrl(context.Background(), code, req)
	if err != nil {
		log.Printf("[UpdateLongUrl] ERROR code=%s error=%s", code, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[UpdateLongUrl] RESPONSE SUCCESS code=%s response=%+v", code, res)
	c.JSON(http.StatusOK, gin.H{"message": res})
}
