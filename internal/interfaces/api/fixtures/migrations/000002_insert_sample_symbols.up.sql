INSERT INTO stock_prices (
    symbol_id, 
    time, 
    open_price, 
    high_price, 
    low_price, 
    close_price, 
    adj_close, 
    volume
)
SELECT 
    s.id,
    price_data.time,
    price_data.open_price,
    price_data.high_price,
    price_data.low_price,
    price_data.close_price,
    price_data.adj_close,
    price_data.volume
FROM symbols s,
(VALUES
    -- Recent trading days with realistic AAPL price data (Eastern Time)
    ('2025-08-07 09:30:00-04'::timestamptz, 228.95, 232.10, 227.40, 231.25, 231.25, 53740000),
    ('2025-08-06 09:30:00-04'::timestamptz, 227.45, 230.20, 225.80, 228.90, 228.90, 49150000),
    ('2025-08-05 09:30:00-04'::timestamptz, 225.85, 229.15, 223.75, 227.60, 227.60, 51280000),
    ('2025-08-02 09:30:00-04'::timestamptz, 224.20, 227.40, 222.10, 225.90, 225.90, 48320000),
    ('2025-08-01 09:30:00-04'::timestamptz, 220.50, 225.80, 218.30, 224.15, 224.15, 52450000),
    ('2025-07-31 09:30:00-04'::timestamptz, 218.80, 222.45, 217.20, 220.50, 220.50, 54210000),
    ('2025-07-30 09:30:00-04'::timestamptz, 216.30, 219.85, 215.10, 218.75, 218.75, 47890000),
    ('2025-07-29 09:30:00-04'::timestamptz, 214.90, 217.60, 213.40, 216.25, 216.25, 50330000),
    ('2025-07-26 09:30:00-04'::timestamptz, 212.15, 215.80, 211.30, 214.85, 214.85, 45670000),
    ('2025-07-25 09:30:00-04'::timestamptz, 210.40, 213.25, 209.20, 212.10, 212.10, 48920000),
    ('2025-07-24 09:30:00-04'::timestamptz, 208.75, 211.50, 207.90, 210.35, 210.35, 46780000),
    ('2025-07-23 09:30:00-04'::timestamptz, 206.20, 209.80, 205.60, 208.70, 208.70, 49240000),
    ('2025-07-22 09:30:00-04'::timestamptz, 204.90, 207.30, 203.80, 206.15, 206.15, 51560000),
    ('2025-07-19 09:30:00-04'::timestamptz, 202.40, 205.95, 201.70, 204.85, 204.85, 48930000),
    ('2025-07-18 09:30:00-04'::timestamptz, 200.80, 203.50, 199.90, 202.35, 202.35, 52180000)
) AS price_data(time, open_price, high_price, low_price, close_price, adj_close, volume)
WHERE s.symbol = 'AAPL'
