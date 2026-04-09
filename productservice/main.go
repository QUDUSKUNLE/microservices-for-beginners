package main

import (
	"log"
	"net/http"
	"os"
	"productservice/internal/db"
	"productservice/internal/handler"
	"productservice/internal/repository"
	"productservice/internal/service"

	"github.com/go-chi/chi"
)

func main() {

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewProductRepo(database)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Post("/{id}/reduce", h.ReduceStock)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from product service"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("Product Service running on ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
