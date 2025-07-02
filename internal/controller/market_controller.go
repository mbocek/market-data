package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"

	"github.com/market-data/pkg/marketservice"
)

// MarketController handles market data related endpoints
type MarketController struct {
	service *marketservice.Service
}

// NewMarketController creates a new market controller
func NewMarketController(service *marketservice.Service) *MarketController {
	return &MarketController{
		service: service,
	}
}

// RegisterRoutes registers the routes for the market controller
func (c *MarketController) RegisterRoutes(router *gin.Engine) {
	router.GET("/symbols", c.getSymbols)
	router.GET("/data/:symbol", c.getMarketData)
}

// getSymbols handles the request to get all available symbols
func (c *MarketController) getSymbols(ctx *gin.Context) {
	symbols := c.service.GetAllSymbols()
	ctx.JSON(http.StatusOK, symbols)
}

// getMarketData handles the request to get market data for a specific symbol
func (c *MarketController) getMarketData(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	marketData, err := c.service.GetMarketData(symbol)
	if err != nil {
		if eris.Is(err, marketservice.ErrSymbolNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Symbol not found"})
		} else {
			log.Error().Err(err).Str("symbol", symbol).Msg("Error getting market data")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, marketData)
}
