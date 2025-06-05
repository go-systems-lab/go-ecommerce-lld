-- Revert orders table back to CHAR(27) for UUID fields
ALTER TABLE orders ALTER COLUMN id TYPE CHAR(27);
ALTER TABLE orders ALTER COLUMN account_id TYPE CHAR(27);

-- Revert order_products table back to CHAR(27) for UUID fields
ALTER TABLE order_products ALTER COLUMN order_id TYPE CHAR(27);
ALTER TABLE order_products ALTER COLUMN product_id TYPE CHAR(27);
