package main

import (
	"authservice/internal/handler"
	"authservice/internal/repository"
	"authservice/internal/service"
	"authservice/internal/telemetry"
	"authservice/pkg/db"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {

	shutdown := telemetry.Init("auth-service")
	defer shutdown()
	db, err := db.InitDB("auth.db")
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewUserRepo(db)
	svc := service.NewAuthService(repo)
	h := handler.NewAuthHandler(svc)
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Print("health hit")
		w.Write([]byte("hello from auth"))
	})
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		otelhttp.NewHandler(http.HandlerFunc(h.Register), "register").ServeHTTP(w, r)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		otelhttp.NewHandler(http.HandlerFunc(h.Login), "login").ServeHTTP(w, r)
	})

	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		otelhttp.NewHandler(http.HandlerFunc(h.Logout), "logout").ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Println("Auth Service running on ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
