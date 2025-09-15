-- orders table (header)
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    order_date TIMESTAMP NOT NULL DEFAULT now(),
    status_staus VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(12,2) NOT NULL,
    payment_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- order_items table (line items)
CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    subtotal NUMERIC(12,2) GENERATED ALWAYS AS (price * quantity) STORED
);