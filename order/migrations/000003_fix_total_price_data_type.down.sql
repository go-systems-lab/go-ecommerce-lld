-- Revert total_price from DECIMAL(10,2) back to MONEY
ALTER TABLE orders ALTER COLUMN total_price TYPE MONEY USING total_price::money;
