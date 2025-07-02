package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMarketData(t *testing.T) {
	// Test case inputs
	symbol := "AAPL"
	price := 175.25
	volume := int64(1000000)

	// Create a new MarketData instance
	md := NewMarketData(symbol, price, volume)

	// Verify the fields were set correctly
	assert.Equal(t, symbol, md.Symbol, "Symbol should match")
	assert.Equal(t, price, md.Price, "Price should match")
	assert.Equal(t, volume, md.Volume, "Volume should match")

	// Verify timestamp is set and is recent
	now := time.Now()
	timeDiff := now.Sub(md.Timestamp)
	assert.LessOrEqual(t, timeDiff, time.Second*5, "Timestamp should be recent")
}

func TestIsValid(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		md       *MarketData
		expected bool
	}{
		{
			name: "Valid market data",
			md: &MarketData{
				Symbol:    "AAPL",
				Price:     175.25,
				Volume:    1000000,
				Timestamp: time.Now(),
			},
			expected: true,
		},
		{
			name: "Empty symbol",
			md: &MarketData{
				Symbol:    "",
				Price:     175.25,
				Volume:    1000000,
				Timestamp: time.Now(),
			},
			expected: false,
		},
		{
			name: "Zero price",
			md: &MarketData{
				Symbol:    "AAPL",
				Price:     0,
				Volume:    1000000,
				Timestamp: time.Now(),
			},
			expected: false,
		},
		{
			name: "Negative price",
			md: &MarketData{
				Symbol:    "AAPL",
				Price:     -10.5,
				Volume:    1000000,
				Timestamp: time.Now(),
			},
			expected: false,
		},
		{
			name: "Negative volume",
			md: &MarketData{
				Symbol:    "AAPL",
				Price:     175.25,
				Volume:    -1000,
				Timestamp: time.Now(),
			},
			expected: false,
		},
		{
			name: "Zero volume is valid",
			md: &MarketData{
				Symbol:    "AAPL",
				Price:     175.25,
				Volume:    0,
				Timestamp: time.Now(),
			},
			expected: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.md.IsValid()
			assert.Equal(t, tc.expected, result, "IsValid() should return expected result")
		})
	}
}
