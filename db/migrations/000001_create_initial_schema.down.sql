-- Drop the price_fetch_logs table
DROP TABLE IF EXISTS price_fetch_logs;

-- Remove the retention policy from stock_prices
-- Note: In TimescaleDB, retention policies are automatically removed when the table is dropped
-- This is included for completeness, but may not be necessary

-- Drop the index on stock_prices
DROP INDEX IF EXISTS idx_stock_prices_symbol_time;

-- Drop the stock_prices table
DROP TABLE IF EXISTS stock_prices;

-- Drop the symbols table
DROP TABLE IF EXISTS symbols;

-- Disable the TimescaleDB extension
DROP EXTENSION IF EXISTS timescaledb;
