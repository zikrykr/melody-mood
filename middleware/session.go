package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melody-mood/pkg"
	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func SessionMiddleware(rds *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Request.Header.Get("X-Session-ID")
		ctx := c.Request.Context()
		if sessionID == "" {
			pkg.ResponseError(c, http.StatusUnauthorized, fmt.Errorf("Missing X-Session-ID header"))
			return
		}

		// Check if the session ID is valid
		valid, err := pkg.CheckValidSession(ctx, rds, sessionID)
		if err != nil {
			logrus.WithError(err).Error("error checking session validity")
			pkg.ResponseError(c, http.StatusInternalServerError, err)
			return
		}

		if !valid {
			pkg.ResponseError(c, http.StatusUnauthorized, fmt.Errorf("Invalid session ID"))
			return
		}

		c.Next()
	}
}
