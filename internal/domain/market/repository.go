package market

import (
	"context"
	"errors"
	"github.com/market-data/internal/providers/yahoo"
	"github.com/rotisserie/eris"

	"github.com/jackc/pgx/v5"
	"github.com/market-data/internal/database"
	"time"
)

// Repository defines the interface for managing financial symbols in a data store.
type Repository interface {
	GetSymbol(ctx context.Context, symbol string) (*Symbol, error)
	GetStockPrice(ctx context.Context, symbol string) (*StockPrices, error)
	SaveSymbol(ctx context.Context, s *Symbol) error
	SaveMarketData(ctx context.Context, data *yahoo.MarketData) error
	SavePriceFetchLogs(ctx context.Context, symbol string, fetchedAt time.Time, dataPoints int,
		success bool, msg string) error
	GetLastFetchTime(ctx context.Context, symbol string) (*time.Time, error)
}

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

// GetSymbol retrieves the details of a specific symbol from the database by its identifier. Returns ErrSymbolNotFound if not found.
func (r *MarketRepository) GetSymbol(ctx context.Context, symbol string) (*Symbol, error) {
	query := `
		SELECT id, symbol, name, exchange, created_at
		FROM symbols
		WHERE symbol = $1
	`

	rows, err := r.db.QueryContext(ctx, query, symbol)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to query symbol: %s", symbol)
	}

	s, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Symbol])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSymbolNotFound
		}
		return nil, eris.Wrapf(err, "cannot collect exactly one row for symbol: %s", symbol)
	}

	return &s, nil
}

// SaveSymbol inserts or updates a symbol record in the symbols table.
func (r *MarketRepository) SaveSymbol(ctx context.Context, s *Symbol) error {
	if err := s.IsValid(); err != nil {
		return eris.Wrapf(err, "invalid symbol: %s", s.Symbol)
	}

	query := `
		INSERT INTO symbols (symbol, name, exchange, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (symbol) DO UPDATE
		SET name = EXCLUDED.name,
			exchange = EXCLUDED.exchange
		RETURNING id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, s.Symbol, s.Name, s.Exchange, time.Now()).Scan(&id)
	if err != nil {
		return eris.Wrapf(err, "failed to save symbol: %s", s.Symbol)
	}
	s.ID = id // Optionally update the symbol struct with the DB id
	return nil
}

func (r *MarketRepository) SaveMarketData(ctx context.Context, data *yahoo.MarketData) error {
	querySymbol := `
		INSERT INTO symbols (symbol, name, exchange, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (symbol) DO UPDATE
		SET name = EXCLUDED.name,
			exchange = EXCLUDED.exchange
		RETURNING id
	`

	queryStockPrice := `
		INSERT INTO stock_prices (
			time,
			symbol_id,
			open_price,
			high_price,
			low_price,
			close_price,
			adj_close,
			volume
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (symbol_id, time) DO UPDATE SET
			open_price = EXCLUDED.open_price,
			high_price = EXCLUDED.high_price,
			low_price = EXCLUDED.low_price,
			close_price = EXCLUDED.close_price,
			adj_close = EXCLUDED.adj_close,
			volume = EXCLUDED.volume;
	`

	return r.db.InTransaction(ctx, func(tx pgx.Tx) error {
		var symbolId int
		err := tx.QueryRow(ctx, querySymbol, data.Symbol, data.Name, data.Exchange, time.Now()).Scan(&symbolId)
		if err != nil {
			return eris.Wrap(err, "failed to insert symbol")
		}

		batch := &pgx.Batch{}
		for _, price := range data.Prices {
			batch.Queue(queryStockPrice,
				price.Time, symbolId, price.Open, price.High, price.Low, price.Close, price.AdjClose, price.Volume)
		}

		// Send the batch and get results
		results := tx.SendBatch(ctx, batch)

		// Always close the batch results when done
		defer results.Close()

		return nil
	})
}

func (r *MarketRepository) SavePriceFetchLogs(ctx context.Context, symbol string, fetchedAt time.Time, dataPoints int,
	success bool, msg string) error {

	queryPriceFetchLog := `
		INSERT INTO price_fetch_logs (
			symbol,     -- nullable INTEGER, FK to symbols(id)
			fetched_at,    -- TIMESTAMPTZ, NOT NULL, defaults to NOW() if not provided
			data_points,   -- nullable INTEGER
			success,       -- BOOLEAN, NOT NULL
			error_msg      -- nullable TEXT
		) VALUES (
			$1,  -- symbol_id
			$2,  -- fetched_at
			$3,  -- data_points
			$4,  -- success
			$5   -- error_msg
		)
		RETURNING id;
	`

	var fetchID int
	err := r.db.QueryRowContext(ctx, queryPriceFetchLog, symbol, fetchedAt, dataPoints, success, msg).Scan(&fetchID)
	if err != nil {
		return eris.Wrap(err, "failed to insert price fetch log")
	}
	return nil
}

func (r *MarketRepository) GetLastFetchTime(ctx context.Context, symbol string) (*time.Time, error) {
	query := `
		SELECT fetched_at
		FROM price_fetch_logs
		WHERE symbol = $1 AND success = true
		ORDER BY fetched_at DESC
		LIMIT 1
	`

	var fetchTime time.Time
	err := r.db.QueryRowContext(ctx, query, symbol).Scan(&fetchTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, eris.Wrapf(err, "failed to get last fetch time for symbol: %s", symbol)
	}

	return &fetchTime, nil
}

// GetStockPrice retrieves stock price data for a specific symbol from the database
func (r *MarketRepository) GetStockPrice(ctx context.Context, symbol string) (*StockPrices, error) {
	query := `
		SELECT sp.time, sp.open_price, sp.high_price, sp.low_price, 
		       sp.close_price, sp.adj_close, sp.volume, sp.symbol_id
		FROM stock_prices sp
		JOIN symbols s ON sp.symbol_id = s.id
		WHERE s.symbol = $1
		ORDER BY sp.time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, symbol)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to query stock prices for symbol: %s", symbol)
	}

	prices, err := pgx.CollectRows(rows, pgx.RowToStructByName[StockPrice])
	if err != nil {
		return nil, eris.Wrapf(err, "failed to collect stock price rows for symbol: %s", symbol)
	}

	if len(prices) == 0 {
		return nil, ErrSymbolNotFound
	}

	stockPrices := StockPrices(prices)
	return &stockPrices, nil
}
