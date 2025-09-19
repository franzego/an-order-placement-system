

INSERT INTO products (product_name, product_description, available_quantity, reserved_quantity, price) VALUES
('Laptop Pro 15"', 'High-performance laptop with 16GB RAM and 512GB SSD', 10.00, 0.00, 1299.99),
('Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 50.00, 0.00, 29.99),
('Mechanical Keyboard', 'RGB mechanical keyboard with Cherry MX switches', 25.00, 0.00, 149.99),
('Monitor 27" 4K', 'Ultra HD 27-inch monitor with HDR support', 8.00, 0.00, 399.99),
('Gaming Headset', '7.1 surround sound gaming headset with microphone', 15.00, 0.00, 89.99),
('USB-C Hub', 'Multi-port USB-C hub with HDMI, USB, and SD card slots', 30.00, 0.00, 49.99),
('External SSD 1TB', 'Portable SSD with USB 3.2 Gen 2 interface', 20.00, 0.00, 199.99),
('Webcam HD', '1080p webcam with autofocus and noise reduction', 12.00, 0.00, 79.99);

SELECT id, product_name, available_quantity, reserved_quantity, price 
FROM products 
ORDER BY id;
