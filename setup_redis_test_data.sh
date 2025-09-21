#!/bin/bash

# Script to populate Redis with test user data
echo "Setting up Redis test data..."

# Connect to Redis and set some test user emails
redis-cli -h localhost -p 6379 SET "123" "testuser123@example.com"
redis-cli -h localhost -p 6379 SET "456" "testuser456@example.com"
redis-cli -h localhost -p 6379 SET "789" "testuser789@example.com"

echo "Test data set:"
echo "User 123 -> testuser123@example.com"
echo "User 456 -> testuser456@example.com"
echo "User 789 -> testuser789@example.com"

# Verify the data
echo -e "\nVerifying data in Redis:"
redis-cli -h localhost -p 6379 GET "123"
redis-cli -h localhost -p 6379 GET "456"
redis-cli -h localhost -p 6379 GET "789"

echo -e "\nAll Redis keys:"
redis-cli -h localhost -p 6379 KEYS "*"
