package repository

import (
	"database/sql"
	"productservice/internal/model"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(p *model.Product) error {
	_, err := r.db.Exec("INSERT INTO products (name,price,category,stock) VALUES(?,?,?,?)", p.Name, p.Price, p.Category, p.Stock)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name,price, category, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepo) GetByID(id int64) (*model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(
		"SELECT id ,name , price, category, stock FROM products WHERE id =?", id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Stock)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) UpdateStock(id int64, stock int) error {
	_, err := r.db.Exec(
		"UPDATE products SET stock = ? WHERE id = ?",
		stock, id,
	)
	return err
}
