package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/harisnkr/url-shortener/core"
	"github.com/joho/godotenv"
)

func setupRoutes() {
	coreAPI := core.New(initRedis())

	router := gin.Default()
	router.GET("/:url", coreAPI.ResolveURL)
	router.POST("/v1/shorten", coreAPI.ShortenURL)

	router.Run(os.Getenv("APP_PORT"))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	setupRoutes()
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       0,
	})
}
