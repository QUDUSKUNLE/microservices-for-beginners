package handler

import (
	"authservice/internal/service"
	"authservice/internal/telemetry"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

type req struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.Tracer().Start(r.Context(), "handler.register")
	defer span.End()

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	err := h.svc.Register(ctx, body.Email, body.Password, body.Address)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created successfully"))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.Tracer().Start(r.Context(), "handler.login")
	defer span.End()

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	token, err := h.svc.Login(ctx, body.Email, body.Password)
	if err != nil {
		http.Error(w, "invalid credentials: "+err.Error(), 401)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logged out"))
}
