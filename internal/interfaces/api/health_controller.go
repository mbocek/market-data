package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check and status endpoints
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// RegisterRoutes registers the routes for the health controller
func (c *HealthController) RegisterRoutes(router *gin.Engine) {
	router.GET("/", c.getStatus)
	router.GET("/health", c.getHealth)
}

// getStatus handles the root endpoint
func (c *HealthController) getStatus(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Market Data Service is running!")
}

// getHealth handles the health check endpoint
func (c *HealthController) getHealth(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}
