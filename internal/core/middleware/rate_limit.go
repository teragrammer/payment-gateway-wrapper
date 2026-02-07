package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	redis2 "github.com/teragrammer/payment-gateway-wrapper/internal/database/redis"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

type RateLimiter struct {
	client      redis.UniversalClient
	maxRequests int
	interval    time.Duration
}

func RateLimitIP() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := utils.GetClientIP(r)
			userAgent := r.UserAgent()
			key := utils.SHA256(ip + ":" + userAgent)

			// Create a throttled Redis store
			redisClient := redis2.ConnectRedis()

			// Set up rate limiter with NUM requests per minute (NUM seconds interval)
			cfg := config.Load()
			rateLimiter := NewRateLimiter(redisClient, cfg.RateLimitIPReq, time.Minute*time.Duration(cfg.RateLimitIPInt))

			ctx := context.Background()
			limited, err := rateLimiter.isRateLimited(ctx, key)
			if err != nil {
				message := fmt.Sprintf("Error checking ip rate limit: %v", err)
				utils.JSONErrorMessage(w, http.StatusUnauthorized, "RATE_LIMIT_IP_CHECK", message)
				return
			}

			if limited {
				message := fmt.Sprintf("Rate limit exceeded for ip address %s", ip)
				utils.JSONErrorMessage(w, http.StatusUnauthorized, "RATE_LIMIT_IP_EXCEED", message)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RateLimitAPIKey() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				utils.JSONErrorMessage(w, http.StatusUnauthorized, "RATE_LIMIT_HEADER", "X-API-Key header not present")
				return
			}

			// Create a throttled Redis store
			redisClient := redis2.ConnectRedis()

			// TODO
			// check if api key exists on projects (mongo)
			// then save to redis for faster fetching

			// Set up rate limiter with NUM requests per minute (NUM seconds interval)
			cfg := config.Load()
			rateLimiter := NewRateLimiter(redisClient, cfg.RateLimitAPIKeyReq, time.Minute*time.Duration(cfg.RateLimitAPIKeyInt))

			ctx := context.Background()
			limited, err := rateLimiter.isRateLimited(ctx, apiKey)
			if err != nil {
				message := fmt.Sprintf("Error checking api key rate limit: %v", err)
				utils.JSONErrorMessage(w, http.StatusUnauthorized, "RATE_LIMIT_API_CHECK", message)
				return
			}

			if limited {
				message := fmt.Sprintf("Rate limit exceeded for api key %s", apiKey)
				utils.JSONErrorMessage(w, http.StatusUnauthorized, "RATE_LIMIT_API_EXCEED", message)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func NewRateLimiter(client redis.UniversalClient, maxRequests int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		client:      client,
		maxRequests: maxRequests,
		interval:    interval,
	}
}

func (rl *RateLimiter) isRateLimited(ctx context.Context, key string) (bool, error) {
	// Redis Key for the rate limit counter
	rateLimitKey := fmt.Sprintf("throttled:%s", key)

	// Use Redis' INCR command to increment the request count
	// We use the current time (as a timestamp) to create the key and apply TTL.
	now := time.Now().Unix()

	// Create a Redis pipeline to optimize multiple operations
	pipe := rl.client.TxPipeline()

	// Increment the current request count in Redis, and set expiration
	count := pipe.Incr(ctx, rateLimitKey)
	pipe.ExpireAt(ctx, rateLimitKey, time.Unix(now+int64(rl.interval.Seconds()), 0))

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Check if the count exceeds maxRequests
	if count.Val() > int64(rl.maxRequests) {
		return true, nil // Rate-limited
	}
	return false, nil
}
