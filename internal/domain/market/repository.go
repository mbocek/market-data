package market

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/market-data/internal/
	"github.com/jackc/pgx/v5"
	"github.com/market-data/internal/database"
	"github.com/rs/zerolog/log"
	"time"
)

// MarketRepository implements the market.Repository interface using PostgreSQL
type MarketRepository struct {
	db *database.DB
}

// NewMarketRepository creates a new market data repository
func NewMarketRepository(db *database.DB) *MarketRepository {
	return &MarketRepository{
		db: db,
	}
}

// GetMarketData retrieves market data for a specific symbol from the database
func (r *MarketRepository) GetMarketData(ctx context.Context, symbol string) (*MarketData, error) {
	query := `
		SELECT s.symbol, sp.close_price, sp.volume, sp.time
		FROM stock_prices sp
		JOIN symbols s ON sp.symbol_id = s.id
		WHERE s.symbol = $1
		ORDER BY sp.time DESC
		LIMIT 1
	`

	var md MarketData
	err := r.db.QueryRowContext(ctx, query, symbol).Scan(
		&md.Symbol,
		&md.Price,
		&md.Volume,
		&md.Timestamp,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrSymbolNotFound
		}
		return nil, eris.Wrapf(err, "failed to get market data for symbol: %s", symbol)
	}

	return &md, nil
}

// SaveMarketData saves market data to the database
func (r *MarketRepository) SaveMarketData(ctx context.Context, md *MarketData) error {
	// Start a transaction
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return eris.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error().Err(rbErr).Msg("Failed to rollback transaction")
			}
		}
	}()

	// Get or create symbol
	var symbolID int
	err = tx.QueryRow(ctx, `
		INSERT INTO symbols (symbol, created_at)
		VALUES ($1, $2)
		ON CONFLICT (symbol) DO UPDATE SET symbol = $1
		RETURNING id
	`, md.Symbol, time.Now()).Scan(&symbolID)

	if err != nil {
		return eris.Wrapf(err, "failed to insert symbol: %s", md.Symbol)
	}

	// Insert stock price
	_, err = tx.Exec(ctx, `
		INSERT INTO stock_prices (time, symbol_id, close_price, volume)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (symbol_id, time) DO UPDATE SET
			close_price = $3,
			volume = $4
	`, md.Timestamp, symbolID, md.Price, md.Volume)

	if err != nil {
		return eris.Wrapf(err, "failed to insert stock price for symbol: %s", md.Symbol)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return eris.Wrap(err, "failed to commit transaction")
	}

	log.Debug().
		Str("symbol", md.Symbol).
		Float64("price", md.Price).
		Int64("volume", md.Volume).
		Time("timestamp", md.Timestamp).
		Msg("Market data saved to database")

	return nil
}

// GetAllSymbols returns a list of all available symbols from the database
func (r *MarketRepository) GetAllSymbols(ctx context.Context) ([]string, error) {
	query := `SELECT symbol FROM symbols ORDER BY symbol`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, eris.Wrap(err, "failed to query symbols")
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return nil, eris.Wrap(err, "failed to scan symbol")
		}
		symbols = append(symbols, symbol)
	}

	if err := rows.Err(); err != nil {
		return nil, eris.Wrap(err, "error iterating symbols")
	}

	return symbols, nil
}
