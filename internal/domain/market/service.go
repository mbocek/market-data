package market

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
)

// Domain errors
var (
	ErrSymbolNotFound = errors.New("symbol not found")
	ErrInvalidData    = errors.New("invalid market data")
)

// MarketService provides core domain operations for market data
type MarketService struct {
	// Repository for persistence
	repo Repository
}

// GetRepository returns the repository used by the service
func (s *MarketService) GetRepository() Repository {
	return s.repo
}

// Repository defines the interface for market data persistence
type Repository interface {
	// GetMarketData retrieves market data for a specific symbol
	GetMarketData(ctx context.Context, symbol string) (*MarketData, error)

	// SaveMarketData saves market data
	SaveMarketData(ctx context.Context, md *MarketData) error
}

// NewMarketService creates a new market data service
func NewMarketService(repo Repository) *MarketService {
	return &MarketService{
		repo: repo,
	}
}

// GetMarketData retrieves market data for a specific symbol
func (s *MarketService) GetMarketData(ctx context.Context, symbol string) (*MarketData, error) {
	log.Debug().Str("symbol", symbol).Msg("Getting market data")
	return s.repo.GetMarketData(ctx, symbol)
}

// UpdateMarketData updates or adds market data for a specific symbol
func (s *MarketService) UpdateMarketData(ctx context.Context, md *MarketData) error {
	if !md.IsValid() {
		return ErrInvalidData
	}

	log.Debug().
		Str("symbol", md.Symbol).
		Float64("price", md.Price).
		Int64("volume", md.Volume).
		Msg("Updating market data")

	return s.repo.SaveMarketData(ctx, md)
}
