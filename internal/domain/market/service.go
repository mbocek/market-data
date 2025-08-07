package market

import (
	"context"
	"errors"
	"github.com/market-data/internal/providers/yahoo"
	"github.com/rotisserie/eris"
	"time"

	"github.com/rs/zerolog/log"
)

// Domain errors
var (
	ErrSymbolNotFound = errors.New("symbol not found")
	ErrInvalidData    = errors.New("invalid market data")
)

// DataProvider defines the interface for market data providers
type DataProvider interface {
	// GetMarketData fetches market data for a specific symbol
	GetMarketData(
		ctx context.Context,
		symbol string,
		interval yahoo.IntervalAPI,
		period yahoo.PeriodAPI,
	) (*yahoo.MarketData, error)
}

// MarketService provides core domain operations for market data
type MarketService struct {
	repo     Repository
	provider DataProvider

	// Auto-update settings
	updateInterval   time.Duration
	enableAutoUpdate bool
	stopChan         chan struct{}

	periodConv   map[int]yahoo.PeriodAPI
	intervalConv map[int]yahoo.IntervalAPI
}

// NewMarketService creates a new market data service
func NewMarketService(repo Repository, provider *yahoo.Client) *MarketService {
	var periodConv = map[int]yahoo.PeriodAPI{
		1:       yahoo.Period1d,
		5:       yahoo.Period5d,
		30:      yahoo.Period1mo,
		3 * 30:  yahoo.Period3mo,
		6 * 30:  yahoo.Period6mo,
		360:     yahoo.Period1y,
		2 * 360: yahoo.Period2y,
		5 * 360: yahoo.Period5y,
	}
	var intervalConv = map[int]yahoo.IntervalAPI{
		1:       yahoo.Interval5m,
		5 * 360: yahoo.Interval1d,
	}

	return &MarketService{
		repo:         repo,
		stopChan:     make(chan struct{}),
		periodConv:   periodConv,
		intervalConv: intervalConv,
		provider:     provider,
	}
}

// SetAutoUpdateSettings sets the auto-update settings for the service
func (s *MarketService) SetAutoUpdateSettings(interval time.Duration, enable bool) {
	s.updateInterval = interval
	s.enableAutoUpdate = enable
}

func (s *MarketService) getPeriodAPI(period int) yahoo.PeriodAPI {
	for i, periodAPI := range s.periodConv {
		if period < i {
			return periodAPI
		}
	}
	return yahoo.Period5y
}

func (s *MarketService) getIntervalAPI(interval int) yahoo.IntervalAPI {
	for i, intervalAPI := range s.intervalConv {
		if interval < i {
			return intervalAPI
		}
	}
	return yahoo.Interval1d
}

//// GetMarketData retrieves market data for a specific symbol
//func (s *MarketService) GetMarketData(ctx context.Context, symbol string) (*MarketData, error) {
//	log.Debug().Str("symbol", symbol).Msg("Getting market data")
//	return s.repo.GetMarketData(ctx, symbol)
//}

//// UpdateMarketData updates or adds market data for a specific symbol
//func (s *MarketService) UpdateMarketData(ctx context.Context, md *MarketData) error {
//	if !md.IsValid() {
//		return ErrInvalidData
//	}
//
//	log.Debug().
//		Str("symbol", md.Symbol).
//		Float64("price", md.Price).
//		Int64("volume", md.Volume).
//		Msg("Updating market data")
//
//	return s.repo.SaveMarketData(ctx, md)
//}

// FetchAndStoreMarketData fetches market data for a symbol from the provider and stores it
func (s *MarketService) FetchAndStoreMarketData(ctx context.Context, symbol string) error {
	if s.provider == nil {
		return errors.New("no data provider configured")
	}

	fetchedAt := time.Now()
	symbolData, err := s.repo.GetSymbol(ctx, symbol)
	if errors.Is(err, ErrSymbolNotFound) {
		data, err := s.provider.GetMarketData(ctx, symbol, yahoo.Interval1d, yahoo.Period5y)
		if err != nil {
			err := s.repo.SavePriceFetchLogs(ctx, symbol, fetchedAt, 0, false, err.Error())
			if err != nil {
				return eris.Wrap(err, "failed to save price fetch logs")
			}
			return eris.Wrap(err, "failed to get market data")
		}

		err = s.repo.SaveMarketData(ctx, data)
		if err != nil {
			err := s.repo.SavePriceFetchLogs(ctx, symbol, fetchedAt, len(data.Prices), false, err.Error())
			if err != nil {
				return eris.Wrap(err, "failed to save price fetch logs")
			}
			return eris.Wrap(err, "failed to save market data")
		}

		err = s.repo.SavePriceFetchLogs(ctx, symbol, fetchedAt, len(data.Prices), true, "")
		if err != nil {
			return eris.Wrap(err, "failed to save price fetch logs")
		}

		return nil
	}

	if err != nil {
		return eris.Wrap(err, "failed to get symbol")
	}

	// get last fetch time for specific symbol form log
	fetchTime, err := s.repo.GetLastFetchTime(ctx, symbol)
	if err != nil {
		return eris.Wrap(err, "failed to get last fetch time")
	}

	if fetchTime == nil {
		return eris.New("fetch time is nil")
	}

	days := int(time.Since(*fetchTime).Hours() / 24)

	intervalAPI := s.getIntervalAPI(days)
	periodAPI := s.getPeriodAPI(days)

	data, err := s.provider.GetMarketData(ctx, symbol, intervalAPI, periodAPI)
	if err != nil {
		err := s.repo.SavePriceFetchLogs(ctx, symbol, fetchedAt, 0, false, err.Error())
		if err != nil {
			return eris.Wrap(err, "failed to save price fetch logs")
		}
		return eris.Wrap(err, "failed to get market data")
	}
	err = s.repo.SaveMarketData(ctx, data)
	if err != nil {
		err := s.repo.SavePriceFetchLogs(ctx, symbol, fetchedAt, len(data.Prices), false, err.Error())
		if err != nil {
			return eris.Wrap(err, "failed to save price fetch logs")
		}
		return eris.Wrap(err, "failed to save market data")
	}

	err = s.repo.SavePriceFetchLogs(ctx, symbol, fetchedAt, len(data.Prices), true, "")
	if err != nil {
		return eris.Wrap(err, "failed to save price fetch logs")
	}

	// download last data
	log.Debug().
		Str("symbol", symbol).
		Interface("symbolData", symbolData).
		Msg("fetching market data from provider")
	return nil
}

func (s *MarketService) GetMarketData(ctx context.Context, symbol string) (*Symbol, *StockPrices, error) {
	symbolData, err := s.repo.GetSymbol(ctx, symbol)
	if err != nil {
		return nil, nil, err
	}

	stockPrices, err := s.repo.GetStockPrice(ctx, symbol)
	if err != nil {
		return nil, nil, err
	}

	return symbolData, stockPrices, nil
}

//// FetchAndStoreAllMarketData fetches market data for all symbols from the provider and stores it
//func (s *MarketService) FetchAndStoreAllMarketData(ctx context.Context) error {
//	if s.provider == nil {
//		return errors.New("no data provider configured")
//	}
//
//	symbols, err := s.repo.GetAllSymbols(ctx)
//	if err != nil {
//		return err
//	}
//
//	if len(symbols) == 0 {
//		log.Info().Msg("No symbols found in database, using default symbols from provider")
//		marketDataList, err := s.provider.GetMarketDataBatch(ctx, nil)
//		if err != nil {
//			return err
//		}
//
//		for _, md := range marketDataList {
//			if err := s.UpdateMarketData(ctx, md); err != nil {
//				log.Error().
//					Err(err).
//					Str("symbol", md.Symbol).
//					Msg("failed to update market data")
//			}
//		}
//
//		return nil
//	}
//
//	log.Info().Int("count", len(symbols)).Msg("fetching market data for all symbols")
//
//	marketDataList, err := s.provider.GetMarketDataBatch(ctx, symbols)
//	if err != nil {
//		return err
//	}
//
//	for _, md := range marketDataList {
//		if err := s.UpdateMarketData(ctx, md); err != nil {
//			log.Error().
//				Err(err).
//				Str("symbol", md.Symbol).
//				Msg("failed to update market data")
//		}
//	}
//
//	return nil
//}
//
//// StartAutoUpdate starts the auto-update process
//func (s *MarketService) StartAutoUpdate() {
//	if !s.enableAutoUpdate || s.provider == nil || s.updateInterval <= 0 {
//		log.Info().Msg("Auto-update not enabled or configured")
//		return
//	}
//
//	log.Info().
//		Dur("interval", s.updateInterval).
//		Msg("Starting auto-update process")
//
//	go func() {
//		ticker := time.NewTicker(s.updateInterval)
//		defer ticker.Stop()
//
//		// Initial update
//		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//		if err := s.FetchAndStoreAllMarketData(ctx); err != nil {
//			log.Error().Err(err).Msg("Failed to fetch market data during initial update")
//		}
//		cancel()
//
//		for {
//			select {
//			case <-ticker.C:
//				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//				if err := s.FetchAndStoreAllMarketData(ctx); err != nil {
//					log.Error().Err(err).Msg("Failed to fetch market data during auto-update")
//				}
//				cancel()
//			case <-s.stopChan:
//				log.Info().Msg("Stopping auto-update process")
//				return
//			}
//		}
//	}()
//}

// StopAutoUpdate stops the auto-update process
func (s *MarketService) StopAutoUpdate() {
	if s.enableAutoUpdate {
		s.stopChan <- struct{}{}
	}
}
