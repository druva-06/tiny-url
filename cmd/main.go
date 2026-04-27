package main

import (
	"log"

	"github.com/druva-06/tiny-url/internal/config"
	"github.com/druva-06/tiny-url/internal/handler"
	"github.com/druva-06/tiny-url/internal/repository"
	"github.com/druva-06/tiny-url/internal/repository/cache"
	"github.com/druva-06/tiny-url/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	db := config.NewDB()
	redis := config.NewRedis()

	rdb := cache.NewURLCache(redis)
	repo := repository.NewURLRepository(db)
	service := service.NewURLService(repo, rdb)
	handler := handler.NewURLHandler(service)

	r := gin.Default()
	r.POST("/url/short", handler.CreateShortURL)
	r.GET("/url/short/:code", handler.GetLongURL)
	r.DELETE("/url/shorten/:code", handler.DeteleShortURL)

	r.Run(":8080")

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		// Fallback in case we are running `go run main.go` from the cmd/ directory
		err = godotenv.Load("../.env")
		if err != nil {
			log.Println("No .env file found (expected in prod)")
		}
	}
}
