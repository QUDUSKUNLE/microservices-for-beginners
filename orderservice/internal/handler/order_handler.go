package handler

import (
	"encoding/json"
	"net/http"
	"orderservice/internal/model"
	"orderservice/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: s}
}

func getEmailFromRequest(r *http.Request) string {
	return r.Header.Get("X-User-Email")
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var o model.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o.UserEmail = getEmailFromRequest(r)

	if o.UserEmail == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.svc.Create(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("order placed successfully"))
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.List(getEmailFromRequest(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
