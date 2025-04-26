package port

import "github.com/gin-gonic/gin"

type ISessionHandler interface {
	// (GET /v1/sessions)
	GenerateSessionID(c *gin.Context)
}
