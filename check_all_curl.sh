#!/bin/bash

# Configuration
API_GATEWAY_URL="http://localhost:8000"

echo "==================================================="
echo " 🚀 Testing Microservices E-commerce API via Gateway"
echo "==================================================="
echo

# 1. Test Auth Service (Login)
echo "[1] Testing Auth Service - Login..."
# Assuming a standard login endpoint. Adjust payload as needed.
LOGIN_RESPONSE=$(curl -s -X POST "$API_GATEWAY_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}')

echo "Response: $LOGIN_RESPONSE"

# Attempt to extract JWT Token 
TOKEN=$(echo $LOGIN_RESPONSE | grep -oE 'eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+' | head -n 1)

if [ -z "$TOKEN" ]; then
    echo "⚠️  Could not automatically extract JWT token from login response."
    echo "You may need to manually paste it below if the Orders route fails."
    TOKEN="your_jwt_token_here"
else
    echo "✅ Token successfully retrieved!"
fi
echo

# 2. Test Product Service (Create Product)
echo "[2] Testing Product Service - Create Product..."
curl -s -X POST "$API_GATEWAY_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name": "Gaming Laptop", "price": 1200.00, "stock": 10}'
echo -e "\n"

# 3. Test Product Service (List Products)
echo "[3] Testing Product Service - List Products..."
curl -s -X GET "$API_GATEWAY_URL/products"
echo -e "\n"

# 4. Test Order Service (Create Order - Protected Route)
echo "[4] Testing Order Service - Create Order (Requires Auth)..."
curl -s -X POST "$API_GATEWAY_URL/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 1}'
echo -e "\n"

# 5. Test Order Service (List Orders - Protected Route)
echo "[5] Testing Order Service - List Orders (Requires Auth)..."
curl -s -X GET "$API_GATEWAY_URL/orders" \
  -H "Authorization: Bearer $TOKEN"
echo -e "\n"

echo "==================================================="
echo " Tests Completed!"
echo "==================================================="
