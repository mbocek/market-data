package api_test

import (
	"encoding/json"
	data "github.com/market-data/db"
	"github.com/market-data/internal/interfaces/api"
	"github.com/market-data/internal/interfaces/api/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/market-data/internal/database"
	"github.com/market-data/internal/domain/market"
	"github.com/market-data/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarketController_GetMarketData(t *testing.T) {
	// Set up TimescaleDB container
	tsdb := tests.NewTimescaleDB(t)
	defer tsdb.Terminate()

	// Apply migrations
	tsdb.ApplyMigrations(data.Migrations)
	tsdb.ApplyMigrationsWitTableName(fixtures.Migrations, "migration_test")

	// Create a new database connection using the host and port from the container
	db, err := database.NewWithConfig(tsdb.DatabaseConfigForTimescale())
	require.NoError(t, err)
	defer db.Close()

	// Initialize repository and service
	repo := market.NewMarketRepository(db)
	service := market.NewMarketService(repo, nil)

	// Insert test data
	testSymbol := "AAPL"
	testPrice := 228.95
	testVolume := int64(53740000)

	//testData := &market.MarketData{
	//	Symbol:    testSymbol,
	//	Price:     testPrice,
	//	Volume:    testVolume,
	//	Timestamp: testTimestamp,
	//}

	//err = service.UpdateMarketData(context.Background(), testData)
	require.NoError(t, err)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Initialize controller and register routes
	controller := api.NewMarketController(service)
	controller.RegisterRoutes(router)

	// Test cases
	t.Run("Get market data for existing symbol", func(t *testing.T) {
		// Create a test request
		req, err := http.NewRequest(http.MethodGet, "/symbols/"+testSymbol, nil)
		require.NoError(t, err)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response body
		var response api.MarketData
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, testSymbol, response.Symbol)
		assert.NotNil(t, response.SymbolPrice)
		assert.Equal(t, testPrice, *response.SymbolPrice[0].Open)
		assert.Equal(t, testVolume, *response.SymbolPrice[0].Volume)
	})

	t.Run("Get market data for non-existent symbol", func(t *testing.T) {
		// Create a test request
		req, err := http.NewRequest(http.MethodGet, "/symbols/NONEXISTENT", nil)
		require.NoError(t, err)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Parse response body
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify error message
		assert.Equal(t, "Symbol not found", response["error"])
	})

	t.Run("Get market data with empty symbol", func(t *testing.T) {
		// Create a test request
		req, err := http.NewRequest(http.MethodGet, "/symbols/", nil)
		require.NoError(t, err)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert response - should be 404 Not Found because the route doesn't match
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
