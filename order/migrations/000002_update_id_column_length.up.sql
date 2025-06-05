-- Update orders table to use VARCHAR(36) for UUID fields
ALTER TABLE orders ALTER COLUMN id TYPE VARCHAR(36);
ALTER TABLE orders ALTER COLUMN account_id TYPE VARCHAR(36);

-- Update order_products table to use VARCHAR(36) for UUID fields  
ALTER TABLE order_products ALTER COLUMN order_id TYPE VARCHAR(36);
ALTER TABLE order_products ALTER COLUMN product_id TYPE VARCHAR(36);
