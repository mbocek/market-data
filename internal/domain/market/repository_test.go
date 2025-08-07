package market_test

import (
	"context"
	_ "embed"
	"encoding/json"
	data "github.com/market-data/db"
	"github.com/market-data/internal/database"
	"github.com/market-data/internal/domain/market"
	"github.com/market-data/internal/domain/market/fixtures"
	"github.com/market-data/internal/providers/yahoo"
	testsTools "github.com/market-data/internal/tests"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	//go:embed fixtures/test/yahoo-data.json
	yahooData []byte
)

func TestMarketRepository_GetSymbol(t *testing.T) {
	timescaleDB := testsTools.NewTimescaleDB(t)
	defer timescaleDB.Terminate()
	timescaleDB.ApplyMigrations(data.Migrations)
	timescaleDB.ApplyMigrationsWitTableName(fixtures.Migrations, "migration_test")

	db, err := database.NewWithConfig(timescaleDB.DatabaseConfigForTimescale())
	require.NoError(t, err)
	defer db.Close()

	marketRepo := market.NewMarketRepository(db)

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{"Get existing symbol", "AAPL", false},
		{"Get non-existent symbol", "NONEXISTENT", true},
		{"Get empty symbol", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := marketRepo.GetSymbol(context.TODO(), tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMarketRepository_SaveSymbol(t *testing.T) {
	timescaleDB := testsTools.NewTimescaleDB(t)
	defer timescaleDB.Terminate()
	timescaleDB.ApplyMigrations(data.Migrations)
	timescaleDB.ApplyMigrationsWitTableName(fixtures.Migrations, "migration_test")

	db, err := database.NewWithConfig(timescaleDB.DatabaseConfigForTimescale())
	require.NoError(t, err)
	defer db.Close()

	marketRepo := market.NewMarketRepository(db)

	tests := []struct {
		name    string
		symbol  *market.Symbol
		wantErr bool
	}{
		{
			name: "Save new symbol",
			symbol: &market.Symbol{
				Symbol:   "GOOG",
				Name:     "Alphabet Inc.",
				Exchange: "NASDAQ",
			},
			wantErr: false,
		},
		{
			name: "Update existing symbol",
			symbol: &market.Symbol{
				Symbol:   "AAPL",
				Name:     "Apple Inc.",
				Exchange: "NASDAQ",
			},
			wantErr: false,
		},
		{
			name: "Save invalid symbol (empty fields)",
			symbol: &market.Symbol{
				Symbol:   "",
				Name:     "",
				Exchange: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := marketRepo.SaveSymbol(context.TODO(), tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMarketRepository_SaveMarketData(t *testing.T) {
	timescaleDB := testsTools.NewTimescaleDB(t)
	defer timescaleDB.Terminate()
	timescaleDB.ApplyMigrations(data.Migrations)

	db, err := database.NewWithConfig(timescaleDB.DatabaseConfigForTimescale())
	require.NoError(t, err)
	defer db.Close()

	marketRepo := market.NewMarketRepository(db)

	var data yahoo.MarketData
	err = json.Unmarshal(yahooData, &data)
	require.NoError(t, err)

	err = marketRepo.SaveMarketData(context.TODO(), &data)
	require.NoError(t, err)
}
