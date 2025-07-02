package data

import (
	"time"
)

// MarketData represents financial market data for a specific instrument
type MarketData struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Volume    int64     `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// NewMarketData creates a new MarketData instance with the current timestamp
func NewMarketData(symbol string, price float64, volume int64) *MarketData {
	return &MarketData{
		Symbol:    symbol,
		Price:     price,
		Volume:    volume,
		Timestamp: time.Now(),
	}
}

// IsValid checks if the market data is valid
func (md *MarketData) IsValid() bool {
	return md.Symbol != "" && md.Price > 0 && md.Volume >= 0
}
