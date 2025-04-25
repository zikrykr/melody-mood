package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/melody-mood/pkg"
	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	sessionLimit = 20            // Max 10 recommendations per session
	sessionTTL   = 1 * time.Hour // Session expires after 1 hour
)

// RateLimitMiddleware enforces a limit of allowed recommendation generations per session
func RateLimitMiddleware(rds *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Request.Header.Get("X-Session-ID")
		ctx := c.Request.Context()
		if sessionID == "" {
			pkg.ResponseError(c, http.StatusUnauthorized, fmt.Errorf("Missing X-Session-ID header"))
			return
		}

		allowed, count, err := allowSessionRequest(ctx, rds, sessionID)
		if err != nil {
			logrus.WithError(err).Error("error checking session rate limit")
			pkg.ResponseError(c, http.StatusInternalServerError, err)
			return
		}

		remaining := sessionLimit - count
		if remaining < 0 {
			remaining = 0
		}

		// Set custom header
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if !allowed {
			pkg.ResponseError(c, http.StatusTooManyRequests, fmt.Errorf("Rate limit exceeded: %d requests allowed per hour", sessionLimit))
			return
		}

		logrus.Infof("Session %s - %d/%d recommendations used", sessionID, count, sessionLimit)
		c.Next()
	}
}

// allowSessionRequest increments the session's request count and checks if it exceeds the limit
func allowSessionRequest(ctx context.Context, rds *redis.Client, sessionID string) (bool, int64, error) {
	key := fmt.Sprintf("session:%s:recommendation_count", sessionID)

	count, err := rds.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	if count == 1 {
		// First request, set expiry
		if err := rds.Expire(ctx, key, sessionTTL).Err(); err != nil {
			return false, count, err
		}
	}

	if count > sessionLimit {
		return false, count, nil
	}

	return true, count, nil
}
