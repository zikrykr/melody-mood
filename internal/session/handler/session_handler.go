package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/session/port"
	"github.com/melody-mood/pkg"
)

type SessionHandler struct {
	sessionService port.ISessionService
}

func NewSessionHandler(service port.ISessionService) port.ISessionHandler {
	return SessionHandler{
		sessionService: service,
	}
}

func (h SessionHandler) GenerateSessionID(c *gin.Context) {
	ctx := c.Request.Context()
	resp, err := h.sessionService.GenerateSessionID(ctx)
	if err != nil {
		pkg.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pkg.HTTPResponse{
		Success: true,
		Message: "Session generated successfully",
		Data:    resp,
	})
}
