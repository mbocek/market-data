package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/market-data/internal/domain/market"
)

type SymbolPrice struct {
	Time     time.Time `json:"time"`
	Open     *float64  `json:"open"`
	High     *float64  `json:"high"`
	Low      *float64  `json:"low"`
	Close    *float64  `json:"close"`
	AdjClose *float64  `json:"adjClose"`
	Volume   *int64    `json:"volume"`
}

type MarketData struct {
	Symbol      string        `json:"symbol"`
	SymbolPrice []SymbolPrice `json:"symbolPrice"`
}

func buildMarketData(symbol *market.Symbol, price *market.StockPrices) *MarketData {
	var symbolPrice []SymbolPrice
	for _, p := range *price {
		symbolPrice = append(symbolPrice, SymbolPrice{
			Time:     p.Time,
			Open:     p.OpenPrice,
			High:     p.HighPrice,
			Low:      p.LowPrice,
			Close:    p.ClosePrice,
			AdjClose: p.AdjClose,
			Volume:   p.Volume,
		})
	}
	return &MarketData{
		Symbol:      symbol.Symbol,
		SymbolPrice: symbolPrice,
	}
}

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

func (c *MarketController) getMarketData(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	symbolData, stockPricesData, err := c.service.GetMarketData(ctx, symbol)
	if err != nil {
		if errors.Is(err, market.ErrSymbolNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Symbol not found"})
			return
		}
		// Handle different types of errors appropriately
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving market data"})
		return
	}

	// Build the response using the helper function
	marketData := buildMarketData(symbolData, stockPricesData)

	ctx.JSON(http.StatusOK, marketData)
}
