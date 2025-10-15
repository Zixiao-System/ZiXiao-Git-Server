#!/bin/bash

# API Test Script for ZiXiao Git Server
# Make sure the server is running before executing this script

set -e

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

echo "======================================"
echo "ZiXiao Git Server - API Test"
echo "======================================"
echo ""

# Test 1: Register a user
echo "[1/7] Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "testpass123",
    "full_name": "Test User"
  }')

if echo "$REGISTER_RESPONSE" | grep -q "token"; then
    TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "✓ User registered successfully"
else
    echo "✗ Registration failed"
    echo "$REGISTER_RESPONSE"
    exit 1
fi

# Test 2: Login
echo "[2/7] Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }')

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    echo "✓ Login successful"
else
    echo "✗ Login failed"
    echo "$LOGIN_RESPONSE"
    exit 1
fi

# Test 3: Get current user
echo "[3/7] Testing get current user..."
USER_RESPONSE=$(curl -s -X GET "$BASE_URL/user" \
  -H "Authorization: Bearer $TOKEN")

if echo "$USER_RESPONSE" | grep -q "testuser"; then
    echo "✓ Get user successful"
else
    echo "✗ Get user failed"
    echo "$USER_RESPONSE"
    exit 1
fi

# Test 4: Create repository
echo "[4/7] Testing repository creation..."
REPO_RESPONSE=$(curl -s -X POST "$BASE_URL/repos" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-repo",
    "description": "Test repository",
    "is_private": false
  }')

if echo "$REPO_RESPONSE" | grep -q "test-repo"; then
    echo "✓ Repository created"
else
    echo "✗ Repository creation failed"
    echo "$REPO_RESPONSE"
    exit 1
fi

# Test 5: Get repository
echo "[5/7] Testing get repository..."
GET_REPO_RESPONSE=$(curl -s -X GET "$BASE_URL/repos/testuser/test-repo")

if echo "$GET_REPO_RESPONSE" | grep -q "test-repo"; then
    echo "✓ Get repository successful"
else
    echo "✗ Get repository failed"
    echo "$GET_REPO_RESPONSE"
    exit 1
fi

# Test 6: List repositories
echo "[6/7] Testing list repositories..."
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/users/testuser/repos")

if echo "$LIST_RESPONSE" | grep -q "test-repo"; then
    echo "✓ List repositories successful"
else
    echo "✗ List repositories failed"
    echo "$LIST_RESPONSE"
    exit 1
fi

# Test 7: Delete repository
echo "[7/7] Testing repository deletion..."
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/repos/testuser/test-repo" \
  -H "Authorization: Bearer $TOKEN")

if echo "$DELETE_RESPONSE" | grep -q "deleted"; then
    echo "✓ Repository deleted"
else
    echo "✗ Repository deletion failed"
    echo "$DELETE_RESPONSE"
    exit 1
fi

echo ""
echo "======================================"
echo "All API tests passed!"
echo "======================================"
echo ""
