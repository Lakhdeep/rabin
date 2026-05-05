#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🧪 Testing Tic-Tac-Toe Authentication API"
echo "========================================"
echo ""

# Test 1: Health Check
echo "1️⃣  Testing Health Check..."
curl -s "$BASE_URL/api/v1/health" | jq '.'
echo ""
echo ""

# Test 2: Register a new user
echo "2️⃣  Testing User Registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "username": "testuser",
    "password": "TestPassword123"
  }')
echo "$REGISTER_RESPONSE" | jq '.'
echo ""
echo ""

# Test 3: Login
echo "3️⃣  Testing User Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "TestPassword123"
  }')
echo "$LOGIN_RESPONSE" | jq '.'

# Extract token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
echo ""
echo "📝 Token: $TOKEN"
echo ""
echo ""

# Test 4: Get Current User
if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    echo "4️⃣  Testing Get Current User (with JWT)..."
    curl -s "$BASE_URL/api/v1/auth/me" \
      -H "Authorization: Bearer $TOKEN" | jq '.'
    echo ""
else
    echo "❌ Login failed - skipping authenticated endpoint test"
fi

echo ""
echo "✅ API tests complete!"
echo ""
echo "To test manually:"
echo "  - Register: curl -X POST $BASE_URL/api/v1/auth/register -H 'Content-Type: application/json' -d '{...}'"
echo "  - Login:    curl -X POST $BASE_URL/api/v1/auth/login -H 'Content-Type: application/json' -d '{...}'"
echo "  - Get User: curl $BASE_URL/api/v1/auth/me -H 'Authorization: Bearer TOKEN'"
