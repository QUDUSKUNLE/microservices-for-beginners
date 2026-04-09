package service

import (
	"errors"
	"productservice/internal/model"
	"productservice/internal/repository"
	"sync"
)

type ProductService struct {
	repo *repository.ProductRepo
	mu   sync.Mutex
}

func NewProductService(r *repository.ProductRepo) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) Create(p *model.Product) error {
	if p.Stock < 0 {
		return errors.New("invalid stock")
	}
	return s.repo.Create(p)
}

func (s *ProductService) List() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Get(id int64) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) ReduceStock(id int64, qty int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if p.Stock < qty {
		return errors.New("not enough stock")
	}

	return s.repo.UpdateStock(id, p.Stock-qty)
}
