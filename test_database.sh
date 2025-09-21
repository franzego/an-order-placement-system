#!/bin/bash

# Test script to verify database connection and data
echo "Testing database connection and data..."

# Check if PostgreSQL is running
echo "1. Checking PostgreSQL connection..."
psql -h localhost -p 5432 -U postgres -d ecommerce -c "SELECT version();"

echo -e "\n2. Checking orders table structure..."
psql -h localhost -p 5432 -U postgres -d ecommerce -c "\d orders"

echo -e "\n3. Checking recent orders..."
psql -h localhost -p 5432 -U postgres -d ecommerce -c "SELECT order_id, user_id, status_staus, total_amount, created_at FROM orders ORDER BY created_at DESC LIMIT 5;"

echo -e "\n4. Checking order_items table..."
psql -h localhost -p 5432 -U postgres -d ecommerce -c "SELECT * FROM order_items ORDER BY order_item_id DESC LIMIT 5;"

echo -e "\nTest completed"
