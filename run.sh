#!/bin/bash

echo "Starting services..."

# run auth service
(cd authservice && go run cmd/main.go) &
AUTH_PID=$!

# run product service
(cd productservice && go run cmd/main.go) &
PRODUCT_PID=$!

# run order service
(cd orderservice && go run cmd/main.go) &
ORDER_PID=$!

# run api gateway
(cd apigateway && go run cmd/main.go) &
GATEWAY_PID=$!

echo "All services started"
echo "Auth PID: $AUTH_PID"
echo "Product PID: $PRODUCT_PID"
echo "Order PID: $ORDER_PID"
echo "Gateway PID: $GATEWAY_PID"

# wait so script doesn't exit
wait