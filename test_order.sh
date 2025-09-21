#!/bin/bash

# Test script to send an order request
echo "Testing order creation with user_id: 123"

curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "args": [
      {
        "product_id": 1001,
        "quantity": 2,
        "price": 4.99
      }
    ]
  }' \
  -v

echo -e "\n\nTest completed"
