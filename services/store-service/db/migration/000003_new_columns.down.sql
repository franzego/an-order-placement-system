ALTER TABLE products
DROP COLUMN IF EXISTS available_quantity,
DROP COLUMN IF EXISTS reserved_quantity;