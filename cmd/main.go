package main

import (
	"fmt"
	"log"
	"os"

	"github.com/druva-06/tiny-url/internal/config"
	"github.com/druva-06/tiny-url/internal/db"
	"github.com/druva-06/tiny-url/internal/handler"
	"github.com/druva-06/tiny-url/internal/repository"
	"github.com/druva-06/tiny-url/internal/repository/cache"
	"github.com/druva-06/tiny-url/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	cluster := db.NewDBCluster()
	redis := config.NewRedis()

	rdb := cache.NewURLCache(redis)
	repo := repository.NewURLRepository(cluster)
	service := service.NewURLService(repo, rdb)
	handler := handler.NewURLHandler(service)

	r := gin.Default()
	fmt.Println("Satrted")
	r.POST("/url/short", handler.CreateShortURL)
	r.GET("/url/short/:code", handler.GetLongURL)
	r.PATCH("/url/short/:code", handler.UpdateLongUrl)

	r.Run(":8080")

}

func loadEnv() {
	// Try both root and cmd-relative locations; missing files are fine in prod/container envs.
	for _, envPath := range []string{".env", "../.env"} {
		if err := godotenv.Load(envPath); err == nil {
			return
		} else if !os.IsNotExist(err) {
			log.Printf("Failed to load %s: %v\n", envPath, err)
		}
	}
}
