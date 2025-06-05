-- Convert total_price from MONEY to DECIMAL(10,2)
-- MONEY type includes currency symbols which can't be parsed as float64
ALTER TABLE orders ALTER COLUMN total_price TYPE DECIMAL(10,2) USING total_price::numeric;
