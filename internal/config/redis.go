package config

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedis() *redis.Client {
	var rdb *redis.Client
	if sentinels := strings.TrimSpace(os.Getenv("REDIS_SENTINELS")); sentinels != "" {
		addrs := strings.Split(sentinels, ",")
		for i := range addrs {
			addrs[i] = strings.TrimSpace(addrs[i])
		}
		rdb = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    os.Getenv("REDIS_MASTER_NAME"),
			SentinelAddrs: addrs,
			Password:      "", // no password by default
			DB:            0,
		})
	} else {
		addr := "localhost:6379"
		if host := strings.TrimSpace(os.Getenv("REDIS_HOST")); host != "" {
			port := strings.TrimSpace(os.Getenv("REDIS_PORT"))
			if port == "" {
				port = "6379"
			}
			addr = host + ":" + port
		}
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password by default
			DB:       0,
		})
	}

	ctx, cancel := context.WithTimeout(Ctx, 2*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err.Error())
	}
	log.Println("Redis connected")
	return rdb
}
