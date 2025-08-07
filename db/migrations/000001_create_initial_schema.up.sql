-- Enable the TimescaleDB extension for time-series support
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- 1. Create the symbols table to store basic information about each ticker
CREATE TABLE IF NOT EXISTS symbols
(
    id         SERIAL PRIMARY KEY,                -- Unique identifier for each symbol
    symbol     TEXT        NOT NULL UNIQUE,       -- Stock ticker symbol (e.g., AAPL)
    name       TEXT        NOT NULL,              -- Full company name
    exchange   TEXT        NOT NULL,              -- Exchange where the symbol is listed (e.g., NASDAQ)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now() -- Timestamp when the record was created
);

-- 2. Create the stock_prices table to store OHLCV data
--    This will later be converted into a hypertable for efficient time-series operations
CREATE TABLE IF NOT EXISTS stock_prices
(
    symbol_id   INTEGER     NOT NULL REFERENCES symbols (id), -- Foreign key to symbols table
    time        TIMESTAMPTZ NOT NULL,                         -- Timestamp of the price data point
    open_price  NUMERIC(18, 6),                               -- Opening price for the interval
    high_price  NUMERIC(18, 6),                               -- Highest price during the interval
    low_price   NUMERIC(18, 6),                               -- Lowest price during the interval
    close_price NUMERIC(18, 6),                               -- Closing price for the interval
    adj_close   NUMERIC(18, 6),                               -- Adjusted closing price accounting for corporate actions
    volume      BIGINT,                                       -- Trading volume during the interval
    PRIMARY KEY (symbol_id, time)                             -- Composite primary key on symbol and time
);

-- Convert the stock_prices table into a TimescaleDB hypertable
--   chunk_time_interval => INTERVAL '1 day' means data is partitioned by day
SELECT create_hypertable(
               'stock_prices',
               'time',
               chunk_time_interval => INTERVAL '1 day',
               if_not_exists => TRUE
       );

-- Create an index to speed up queries filtering by symbol and time in descending order
CREATE INDEX IF NOT EXISTS idx_stock_prices_symbol_time
    ON stock_prices (symbol_id, time DESC);

-- Optional: Add a retention policy to automatically drop data older than 1 year
SELECT add_retention_policy('stock_prices', INTERVAL '1 year');

-- 3. Create the price_fetch_logs table to audit data fetch operations
CREATE TABLE IF NOT EXISTS price_fetch_logs
(
    id          SERIAL PRIMARY KEY,                 -- Unique identifier for each log entry
    symbol      TEXT,                               -- Which symbol was fetched
    fetched_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- Timestamp of when the fetch occurred
    data_points INTEGER,                            -- Number of data points retrieved during fetch
    success     BOOLEAN     NOT NULL,               -- Whether the fetch operation succeeded
    error_msg   TEXT                                -- Error message if the fetch failed
);