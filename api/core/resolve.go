package core

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/harisnkr/url-shortener/errors"
)

func (i *impl) ResolveURL(c *gin.Context) {
	var (
		url = c.Param("url")
		ctx = context.Background()
	)

	// to apply rate limit on resolving URLs

	value, err := i.cache.Get(ctx, url).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, errors.ShortenedURLNotFoundError)
			return
		}

		c.JSON(http.StatusInternalServerError, errors.DatabaseError)
		return
	}

	_ = i.cache.Incr(ctx, "counter")

	c.Redirect(http.StatusMovedPermanently, value)
	return
}
