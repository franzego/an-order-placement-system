ALTER TABLE store
ADD COLUMN available_quantity numeric(10,2) not null,
ADD COLUMN reserved_quantity numeric(10,2) not null;
