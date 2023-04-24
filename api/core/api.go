package core

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func New(cache *redis.Client) ShortenerAPI {
	return &impl{cache: cache}
}

type impl struct {
	cache *redis.Client
}

type ShortenerAPI interface {
	ShortenURL(c *gin.Context)
	ResolveURL(c *gin.Context)
}
