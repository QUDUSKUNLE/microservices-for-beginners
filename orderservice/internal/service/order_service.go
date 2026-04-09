package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"orderservice/internal/model"
	"orderservice/internal/producer"
	"orderservice/internal/repository"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type OrderService struct {
	repo *repository.OrderRepo
	mu   sync.Mutex
}

func NewOrderService(r *repository.OrderRepo) *OrderService {
	return &OrderService{repo: r}
}

type ProductResponse struct {
	ID    int64 `json:"id"`
	Stock int   `json:"stock"`
}

func (s *OrderService) Create(o *model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	productBaseURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productBaseURL == "" {
		productBaseURL = "http://localhost:8002"
	}
	url := productBaseURL + "/" + strconv.FormatInt(o.ProductID, 10)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("product not found")
	}

	var p ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return err
	}

	// ✅ check BEFORE reducing
	if p.Stock < o.Quantity {
		return errors.New("not enough stock")
	}

	// reduce stock
	reduceURL := url + "/reduce"

	body := map[string]int{"quantity": o.Quantity}
	jsonBody, _ := json.Marshal(body)

	resp2, err := http.Post(reduceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		return errors.New("failed to reduce stock")
	}

	// create order
	o.Status = "PLACED"
	if err := s.repo.Create(o); err != nil {
		return err
	}

	// publish event
	oce := producer.OrderCreatedEvent{
		EventID:   uuid.New().String(),
		UserEmail: o.UserEmail,
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
	}

	if err := producer.PublishOrderEvent(&oce); err != nil {
		log.Println("error in producer", err)
	}

	return nil
}	

func (s *OrderService) List(user string) ([]model.Order, error) {
	return s.repo.GetAll(user)
}
