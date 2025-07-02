package marketservice

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"

	"github.com/market-data/internal/data"
	"github.com/market-data/internal/database"
)

// Errors
var (
	ErrSymbolNotFound = eris.New("symbol not found")
	ErrInvalidData    = eris.New("invalid market data")
)

// Service provides methods for retrieving market data
type Service struct {
	// In-memory cache
	cache map[string]*data.MarketData
	mutex sync.RWMutex

	// Database repository
	repo       *data.Repository
	ctxTimeout time.Duration
}

// ServiceOption is a function that configures a Service
type ServiceOption func(*Service)

// WithDatabase configures the service to use a database
func WithDatabase(db *database.DB) ServiceOption {
	return func(s *Service) {
		s.repo = data.NewRepository(db)
	}
}

// WithContextTimeout configures the context timeout for database operations
func WithContextTimeout(timeout time.Duration) ServiceOption {
	return func(s *Service) {
		s.ctxTimeout = timeout
	}
}

// NewService creates a new market data service
func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		cache:      make(map[string]*data.MarketData),
		mutex:      sync.RWMutex{},
		ctxTimeout: 5 * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	log.Debug().Msg("Creating new market data service")
	return s
}

// GetMarketData retrieves market data for a specific symbol
func (s *Service) GetMarketData(symbol string) (*data.MarketData, error) {
	log.Debug().Str("symbol", symbol).Msg("Getting market data")

	// Try to get from database first
	if s.repo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
		defer cancel()

		md, err := s.repo.GetMarketData(ctx, symbol)
		if err == nil {
			// Cache the result
			s.mutex.Lock()
			s.cache[symbol] = md
			s.mutex.Unlock()
			return md, nil
		}

		// If not found in DB, fall back to cache
		if !eris.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Str("symbol", symbol).Msg("Database error getting market data")
		}
	}

	// Get from cache
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	md, ok := s.cache[symbol]
	if !ok {
		return nil, eris.Wrapf(ErrSymbolNotFound, "symbol: %s", symbol)
	}

	return md, nil
}

// UpdateMarketData updates or adds market data for a specific symbol
func (s *Service) UpdateMarketData(md *data.MarketData) error {
	if !md.IsValid() {
		return eris.Wrapf(ErrInvalidData, "symbol: %s, price: %f, volume: %d", 
			md.Symbol, md.Price, md.Volume)
	}

	log.Debug().
		Str("symbol", md.Symbol).
		Float64("price", md.Price).
		Int64("volume", md.Volume).
		Msg("Updating market data")

	// Update cache
	s.mutex.Lock()
	s.cache[md.Symbol] = md
	s.mutex.Unlock()

	// Update database
	if s.repo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
		defer cancel()

		if err := s.repo.SaveMarketData(ctx, md); err != nil {
			log.Error().Err(err).Str("symbol", md.Symbol).Msg("Failed to save market data to database")
			return eris.Wrap(err, "failed to save to database")
		}
	}

	return nil
}

// GetAllSymbols returns a list of all available symbols
func (s *Service) GetAllSymbols() []string {
	log.Debug().Msg("Getting all symbols")

	// Get from database
	if s.repo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.ctxTimeout)
		defer cancel()

		symbols, err := s.repo.GetAllSymbols(ctx)
		if err == nil {
			log.Debug().Int("count", len(symbols)).Msg("Retrieved symbols from database")
			return symbols
		}

		log.Error().Err(err).Msg("Failed to get symbols from database, falling back to cache")
	}

	// Fall back to cache
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	symbols := make([]string, 0, len(s.cache))
	for symbol := range s.cache {
		symbols = append(symbols, symbol)
	}

	log.Debug().Int("count", len(symbols)).Msg("Retrieved symbols from cache")
	return symbols
}
