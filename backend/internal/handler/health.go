package handler

import (
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gin-gonic/gin"
)

// Healthz is the liveness endpoint.
func Healthz(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"})
}

// Health is used by Docker healthcheck.
func Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}
