package repository

import (
	"database/sql"
	"orderservice/internal/model"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(o *model.Order) error {
	_, err := r.db.Exec("INSERT INTO orders (user_email,product_id,quantity,address,status)VALUES(?,?,?,?,?)",
		o.UserEmail, o.ProductID, o.Quantity, o.Address, o.Status)
	return err
}

func (r *OrderRepo) GetAll(user string) ([]model.Order, error) {
	rows, err := r.db.Query(
		"SELECT id, user_email,product_id,quantity,address,status FROM orders WHERE user_email=?",
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserEmail, &o.ProductID, &o.Quantity, &o.Address, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
