package main

import (
	"apigateway/internal/handler"
	"apigateway/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RateLimit)

	// service urls from env (defaults for local dev)
	authService := os.Getenv("AUTH_SERVICE_URL")
	if authService == "" {
		authService = "http://localhost:8001"
	}
	productService := os.Getenv("PRODUCT_SERVICE_URL")
	if productService == "" {
		productService = "http://localhost:8002"
	}
	orderService := os.Getenv("ORDER_SERVICE_URL")
	if orderService == "" {
		orderService = "http://localhost:8003"
	}

	// proxies
	authProxy := handler.NewProxy(authService, "/auth")
	productProxy := handler.NewProxy(productService, "/products")
	orderProxy := handler.NewProxy(orderService, "/orders")
	r.Mount("/auth", authProxy)
	r.Mount("/products", productProxy)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT)
		r.Mount("/orders", orderProxy)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" 
	}

	log.Println("API Gateway running on ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
