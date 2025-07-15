#!/bin/bash

# Test script for Sykell Backend Web Crawler API

BASE_URL="http://localhost:8080/api"

echo "=== Sykell Web Crawler API Test ==="
echo

# Test 1: Create a user and login to get token
echo "1. Creating user and logging in..."
USER_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"testpassword"}')

echo "User created: $USER_RESPONSE"

LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","password":"testpassword"}')

echo "Login response: $LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token: $TOKEN"
echo

# Test 2: Add a URL for crawling
echo "2. Adding URL for crawling..."
ADD_URL_RESPONSE=$(curl -s -X POST $BASE_URL/crawler/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"url":"https://example.com"}')

echo "Add URL response: $ADD_URL_RESPONSE"

# Extract URL ID
URL_ID=$(echo $ADD_URL_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "URL ID: $URL_ID"
echo

# Test 3: Start crawling
echo "3. Starting crawl..."
START_CRAWL_RESPONSE=$(curl -s -X POST $BASE_URL/crawler/urls/$URL_ID/crawl \
  -H "Authorization: Bearer $TOKEN")

echo "Start crawl response: $START_CRAWL_RESPONSE"
echo

# Test 4: Wait a bit and check results
echo "4. Waiting 10 seconds for crawl to complete..."
sleep 10

CRAWL_RESULT=$(curl -s -X GET $BASE_URL/crawler/urls/$URL_ID \
  -H "Authorization: Bearer $TOKEN")

echo "Crawl result: $CRAWL_RESULT"
echo

# Test 5: Get all crawl URLs
echo "5. Getting all crawl URLs..."
ALL_URLS_RESPONSE=$(curl -s -X GET $BASE_URL/crawler/urls \
  -H "Authorization: Bearer $TOKEN")

echo "All URLs response: $ALL_URLS_RESPONSE"
echo

# Test 6: Get crawler stats
echo "6. Getting crawler statistics..."
STATS_RESPONSE=$(curl -s -X GET $BASE_URL/crawler/stats \
  -H "Authorization: Bearer $TOKEN")

echo "Stats response: $STATS_RESPONSE"
echo

# Test 7: Bulk add URLs
echo "7. Bulk adding URLs..."
BULK_ADD_RESPONSE=$(curl -s -X POST $BASE_URL/crawler/urls/bulk \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"urls":["https://httpbin.org","https://github.com","https://stackoverflow.com"]}')

echo "Bulk add response: $BULK_ADD_RESPONSE"
echo

echo "=== Test completed ==="
echo "Note: Make sure MySQL is running and the server is started before running these tests."
