package producer

import (
	"encoding/json"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type OrderCreatedEvent struct {
	EventID   string `json:"event_id"`
	UserEmail string `json:"user_email"`
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func PublishOrderEvent(o *OrderCreatedEvent) error {

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declare DLX exchange
	if err := ch.ExchangeDeclare("dlx", "direct", true, false, false, false, nil); err != nil {
		return err
	}

	// declare DLQ
	if _, err := ch.QueueDeclare("orders-dlq", true, false, false, false, nil); err != nil {
		return err
	}

	// bind DLQ to DLX
	if err := ch.QueueBind("orders-dlq", "", "dlx", false, nil); err != nil {
		return err
	}

	args := amqp091.Table{
		"x-dead-letter-exchange": "dlx",
	}

	q, err := ch.QueueDeclare(
		"orders",
		true, // durable
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(o)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
}
