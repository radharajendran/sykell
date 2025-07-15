#!/bin/bash

# Test script for Sykell Backend API with MySQL

BASE_URL="http://localhost:8080/api"

echo "=== Sykell Backend API Test ==="
echo

# Test 1: Create a user
echo "1. Creating a new user..."
USER_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"testpassword"}')

echo "Response: $USER_RESPONSE"
echo

# Test 2: Login with the created user
echo "2. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","password":"testpassword"}')

echo "Response: $LOGIN_RESPONSE"

# Extract token from response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token: $TOKEN"
echo

# Test 3: Get all users (with authentication)
echo "3. Getting all users (authenticated)..."
USERS_RESPONSE=$(curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $USERS_RESPONSE"
echo

# Test 4: Get user by ID
echo "4. Getting user by ID..."
USER_ID_RESPONSE=$(curl -s -X GET $BASE_URL/users/1 \
  -H "Authorization: Bearer $TOKEN")

echo "Response: $USER_ID_RESPONSE"
echo

echo "=== Test completed ==="
echo "Note: Make sure MySQL is running and the database is set up before running these tests."
