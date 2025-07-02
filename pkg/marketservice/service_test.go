package marketservice

import (
	"testing"

	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/market-data/internal/data"
)

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service, "NewService() should return a non-nil service")
	assert.NotNil(t, service.cache, "service.cache should be initialized")
}

func TestGetMarketData(t *testing.T) {
	// Create a new service
	service := NewService()

	// Test getting non-existent symbol
	_, err := service.GetMarketData("NONEXISTENT")
	assert.Error(t, err, "Getting non-existent symbol should return an error")
	assert.True(t, eris.Is(err, ErrSymbolNotFound), "Error should be ErrSymbolNotFound")

	// Add a symbol
	symbol := "AAPL"
	price := 175.25
	volume := int64(1000000)
	md := data.NewMarketData(symbol, price, volume)

	// Update the market data
	err = service.UpdateMarketData(md)
	require.NoError(t, err, "Updating market data should not return an error")

	// Get the market data
	retrievedMd, err := service.GetMarketData(symbol)
	require.NoError(t, err, "Getting market data should not return an error")

	// Verify the data
	assert.Equal(t, symbol, retrievedMd.Symbol, "Symbol should match")
	assert.Equal(t, price, retrievedMd.Price, "Price should match")
	assert.Equal(t, volume, retrievedMd.Volume, "Volume should match")
}

func TestUpdateMarketData(t *testing.T) {
	// Create a new service
	service := NewService()

	// Test updating with invalid market data
	invalidMd := &data.MarketData{
		Symbol:    "", // Invalid: empty symbol
		Price:     175.25,
		Volume:    1000000,
		Timestamp: data.NewMarketData("TEMP", 0, 0).Timestamp, // Just to get a valid timestamp
	}

	err := service.UpdateMarketData(invalidMd)
	assert.Error(t, err, "Updating invalid market data should return an error")
	assert.True(t, eris.Is(err, ErrInvalidData), "Error should be ErrInvalidData")

	// Test updating with valid market data
	validMd := data.NewMarketData("AAPL", 175.25, 1000000)
	err = service.UpdateMarketData(validMd)
	require.NoError(t, err, "Updating valid market data should not return an error")

	// Verify the data was stored
	retrievedMd, err := service.GetMarketData("AAPL")
	require.NoError(t, err, "Getting market data should not return an error")

	assert.Equal(t, validMd, retrievedMd, "Retrieved market data should match the original")
}

func TestGetAllSymbols(t *testing.T) {
	// Create a new service
	service := NewService()

	// Initially, there should be no symbols
	symbols := service.GetAllSymbols()
	assert.Empty(t, symbols, "Initially, there should be no symbols")

	// Add some market data
	testData := []*data.MarketData{
		data.NewMarketData("AAPL", 175.25, 1000000),
		data.NewMarketData("MSFT", 325.50, 750000),
		data.NewMarketData("GOOGL", 135.75, 500000),
	}

	for _, md := range testData {
		err := service.UpdateMarketData(md)
		require.NoError(t, err, "Updating market data should not return an error")
	}

	// Get all symbols
	symbols = service.GetAllSymbols()

	// Verify the number of symbols
	assert.Len(t, symbols, len(testData), "Number of symbols should match")

	// Verify each symbol is in the list
	expectedSymbols := map[string]bool{
		"AAPL":  true,
		"MSFT":  true,
		"GOOGL": true,
	}

	for _, symbol := range symbols {
		assert.True(t, expectedSymbols[symbol], "Symbol should be in the expected list")
		// Mark as found
		expectedSymbols[symbol] = false
	}

	// Verify all expected symbols were found
	for symbol, notFound := range expectedSymbols {
		assert.False(t, notFound, "Symbol %s should be found", symbol)
	}
}
