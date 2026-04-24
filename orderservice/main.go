package main

import (
	"log"
	"net/http"
	"orderservice/internal/db"
	"orderservice/internal/handler"
	"orderservice/internal/repository"
	"orderservice/internal/service"
	"os"

	"github.com/go-chi/chi"
)

func main() {

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewOrderRepo(database)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from Order Service"))
	})

	r.Post("/", h.Create)
	r.Get("/", h.List)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	log.Println("Order Service running on ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}

