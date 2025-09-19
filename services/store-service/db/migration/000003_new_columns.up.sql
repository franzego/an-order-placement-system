ALTER TABLE products
ADD COLUMN available_quantity numeric(10,2) not null DEFAULT 0,
ADD COLUMN reserved_quantity numeric(10,2) not null DEFAULT 0;
