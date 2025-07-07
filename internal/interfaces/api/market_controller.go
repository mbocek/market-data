package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"

	"github.com/market-data/internal/domain/market"
)

// MarketController handles market data related endpoints
type MarketController struct {
	service *market.MarketService
}

// NewMarketController creates a new market controller
func NewMarketController(service *market.MarketService) *MarketController {
	return &MarketController{
		service: service,
	}
}

// RegisterRoutes registers the routes for the market controller
func (c *MarketController) RegisterRoutes(router *gin.Engine) {
	router.GET("/symbols/:symbol", c.getMarketData)
}

// getMarketData handles the request to get market data for a specific symbol
func (c *MarketController) getMarketData(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	marketData, err := c.service.GetMarketData(ctx, symbol)
	if err != nil {
		if eris.Is(err, market.ErrSymbolNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Symbol not found"})
		} else {
			log.Error().Err(err).Str("symbol", symbol).Msg("Error getting market data")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, marketData)
}
