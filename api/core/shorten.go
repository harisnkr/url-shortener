package core

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/harisnkr/url-shortener/errors"
	"github.com/harisnkr/url-shortener/handlers"
)

type shortenRequest struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	Expiry      time.Duration `json:"expiry"`
}

type shortenResponse struct {
	URL                 string        `json:"url"`
	CustomShort         string        `json:"customShort"`
	Expiry              time.Duration `json:"expiry"`
	XRateLimitRemaining int           `json:"requestsRemaining"`
	XRateLimitResetsIn  int64         `json:"limitResetsIn"`
}

func (i *impl) ShortenURL(c *gin.Context) {
	var (
		ctx = context.Background()
	)

	requestBody := &shortenRequest{}
	if err := c.ShouldBindBodyWith(requestBody, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errors.InvalidParametersError)
	}

	val, err := i.cache.Get(ctx, c.ClientIP()).Result()
	if err != nil {
		if err == redis.Nil {
			i.cache.Set(ctx, c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
		}
		c.JSON(http.StatusInternalServerError, errors.DatabaseError)
		return
	}

	val, _ = i.cache.Get(ctx, c.ClientIP()).Result()
	valInt, _ := strconv.Atoi(val)
	if valInt <= 0 {
		limit, _ := i.cache.TTL(ctx, c.ClientIP()).Result()
		c.JSON(http.StatusTooManyRequests, &errors.ErrorResponse{
			Error:              "Too many requests",
			XRateLimitResetsIn: int64(limit / time.Nanosecond / time.Minute),
		})
		return
	}

	if !govalidator.IsURL(requestBody.URL) {
		c.JSON(http.StatusBadRequest, errors.InvalidURLError)
		return
	}

	if !handlers.IsValidDomain(requestBody.URL) {
		c.JSON(http.StatusServiceUnavailable, errors.InvalidDomainError)
		return
	}

	requestBody.URL = handlers.EnforceHTTP(requestBody.URL)

	var id string
	if requestBody.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = requestBody.CustomShort
	}

	val, _ = i.cache.Get(ctx, id).Result()
	if val != "" {
		c.JSON(http.StatusConflict, errors.URLInUseError)
		return
	}

	if requestBody.Expiry == 0 {
		requestBody.Expiry = 24
	}

	if err = i.cache.Set(ctx, id, requestBody.URL, requestBody.Expiry*3600*time.Second).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.DatabaseError)
		return
	}

	resp := &shortenResponse{
		URL:                 requestBody.URL,
		Expiry:              requestBody.Expiry,
		CustomShort:         "http://" + os.Getenv("DOMAIN") + "/" + id,
		XRateLimitRemaining: 10,
		XRateLimitResetsIn:  30,
	}

	i.cache.Decr(ctx, c.ClientIP())

	val, _ = i.cache.Get(ctx, c.ClientIP()).Result()
	resp.XRateLimitRemaining, _ = strconv.Atoi(val)

	ttl, _ := i.cache.TTL(ctx, c.ClientIP()).Result()
	resp.XRateLimitResetsIn = int64(ttl / time.Nanosecond / time.Minute)

	c.JSON(http.StatusOK, resp)
	return

}
