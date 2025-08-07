package yahoo_test

import (
	"context"
	"github.com/market-data/internal/config"
	"github.com/market-data/internal/providers/yahoo"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestYahoo(t *testing.T) {
	// manual test
	if os.Getenv("MANUAL_TEST") == "" {
		t.Skip("Skipping this test when running from CLI")
	}

	client := yahoo.NewClient(&config.YahooFinanceConfig{
		BaseURL:        "https://query1.finance.yahoo.com/v8/finance/chart/",
		RequestTimeout: 10,
		RetryCount:     3,
	})

	data, err := client.GetMarketData(context.TODO(), "AAPL", yahoo.Interval1d, yahoo.Period1mo)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	log.Info().Interface("data", data).Msg("Data")
}
