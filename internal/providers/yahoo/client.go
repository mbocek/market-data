package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/market-data/internal/config"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

// Client is a Yahoo Finance API client
type Client struct {
	baseURL        string
	httpClient     *http.Client
	retryCount     int
	retryWaitTime  time.Duration
	defaultSymbols []string
}

// GetMarketData retrieves market data for a given symbol, interval, and period from Yahoo Finance API.
func (c *Client) GetMarketData(
	ctx context.Context,
	symbol string,
	interval IntervalAPI,
	period PeriodAPI,
) (*MarketData, error) {
	url := fmt.Sprintf("%s%s?interval=%s&range=%s", c.baseURL, symbol, interval, period)

	var resp *http.Response
	var err error

	// Retry logic
	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			log.Debug().
				Str("symbol", symbol).
				Int("attempt", attempt).
				Msg("Retrying Yahoo Finance API request")

			select {
			case <-ctx.Done():
				return nil, eris.Wrap(ctx.Err(), "context canceled while waiting to retry")
			case <-time.After(c.retryWaitTime):
				// Continue with retry
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, eris.Wrap(err, "failed to create request")
		}

		// Add User-Agent header here
		req.Header.Set("User-Agent", "MarketDataService/1.0")

		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				return nil, eris.Wrap(err, "cannot close response body")
			}
		}

		if err == nil {
			err = eris.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		if attempt == c.retryCount {
			return nil, eris.Wrapf(err, "failed to fetch market data for symbol %s after %d attempts", symbol, c.retryCount+1)
		}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, eris.Wrap(err, "failed to read response body")
	}

	var yahooResp YahooFinanceResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal Yahoo Finance response")
	}

	return c.transform(yahooResp)
}

func (c *Client) transform(resp YahooFinanceResponse) (*MarketData, error) {
	results := resp.Chart.Result
	if len(results) == 0 {
		return nil, eris.New("no results found")
	}
	result := results[0]
	marketData := &MarketData{
		Symbol:   result.Meta.Symbol,
		Name:     result.Meta.Name,
		Exchange: result.Meta.ExchangeName,
	}
	for i, ts := range result.Timestamp {
		marketData.Prices = append(marketData.Prices, StockPrice{
			Time:     time.Unix(ts, 0),
			Open:     result.Indicators.Quote[0].Open[i],
			High:     result.Indicators.Quote[0].High[i],
			Low:      result.Indicators.Quote[0].Low[i],
			Close:    result.Indicators.Quote[0].Close[i],
			AdjClose: result.Indicators.Adjclose[0].Adjclose[i],
			Volume:   int(result.Indicators.Quote[0].Volume[i]),
		})
	}
	return marketData, nil
}

// NewClient creates a new Yahoo Finance client
func NewClient(cfg *config.YahooFinanceConfig) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.GetRequestTimeout(),
		},
		retryCount:     cfg.RetryCount,
		retryWaitTime:  cfg.GetRetryWaitTime(),
		defaultSymbols: cfg.DefaultSymbols,
	}
}
