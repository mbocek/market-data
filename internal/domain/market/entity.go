package market

import (
	"errors"
	"github.com/market-data/internal/providers/yahoo"
	"time"
)

// Validation errors
var (
	ErrSymbolRequired = errors.New("symbol is required")
)

// Symbol represents a financial symbol with its metadata and creation timestamp.
type Symbol struct {
	ID        int       `db:"id"`
	Symbol    string    `db:"symbol"`
	Name      string    `db:"name"`
	Exchange  string    `db:"exchange"`
	CreatedAt time.Time `db:"created_at"`
}

func NewSymbolFromMarketData(m *yahoo.MarketData) *Symbol {
	return &Symbol{
		Symbol:    m.Symbol,
		Name:      m.Name,
		Exchange:  m.Exchange,
		CreatedAt: time.Now(),
	}
}

// IsValid validates the Symbol object and returns an error if the `Symbol` field is empty.
func (s *Symbol) IsValid() error {
	if s.Symbol == "" {
		return ErrSymbolRequired
	}

	return nil
}

type StockPrices []StockPrice

// StockPrice represents the data for a specific stock's price at a particular point in time.
type StockPrice struct {
	Time       time.Time `db:"time"`        // TIMESTAMPTZ NOT NULL
	SymbolID   int       `db:"symbol_id"`   // INTEGER NOT NULL (FK)
	OpenPrice  *float64  `db:"open_price"`  // NUMERIC(18, 6) nullable
	HighPrice  *float64  `db:"high_price"`  // NUMERIC(18, 6) nullable
	LowPrice   *float64  `db:"low_price"`   // NUMERIC(18, 6) nullable
	ClosePrice *float64  `db:"close_price"` // NUMERIC(18, 6) nullable
	AdjClose   *float64  `db:"adj_close"`   // NUMERIC(18, 6) nullable
	Volume     *int64    `db:"volume"`      // BIGINT nullable
}

// PriceFetchLog represents a log of price fetching activity, including metadata such as success status and error details.
type PriceFetchLog struct {
	ID         int       `db:"id"`          // SERIAL PRIMARY KEY
	SymbolID   *int      `db:"symbol_id"`   // INTEGER (nullable), FK to symbols(id)
	FetchedAt  time.Time `db:"fetched_at"`  // TIMESTAMPTZ NOT NULL DEFAULT now()
	DataPoints *int      `db:"data_points"` // INTEGER (nullable)
	Success    bool      `db:"success"`     // BOOLEAN NOT NULL
	ErrorMsg   *string   `db:"error_msg"`   // TEXT (nullable)
}
