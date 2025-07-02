package controller

import (
	"github.com/gin-gonic/gin"
)

// Controller defines the interface for all controllers
type Controller interface {
	// RegisterRoutes registers the routes for this controller
	RegisterRoutes(router *gin.Engine)
}
