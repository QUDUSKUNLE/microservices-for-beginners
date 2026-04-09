package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	_ "modernc.org/sqlite"
)

type OrderCreatedEvent struct {
	EventID   string `json:"event_id"`
	UserEmail string `json:"user_email"`
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type SQliteDB struct {
	db *sql.DB
}

func InitDB(dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS processed_events (
		event_id TEXT PRIMARY KEY
	);`

	_, err = db.Exec(schema)
	return db, err
}

func (s *SQliteDB) processEvent(id string) (bool, error) {
	var existing string

	query := "SELECT event_id FROM processed_events WHERE event_id = ?"
	err := s.db.QueryRow(query, id).Scan(&existing)
	if err == nil {
		log.Println("already processed")
		return false, nil
	}

	insertQuery := "INSERT INTO processed_events (event_id) VALUES(?)"
	if _, err := s.db.Exec(insertQuery, id); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			log.Println("already exists")
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func main() {

	database, err := InitDB("notifications.db")
	if err != nil {
		log.Fatal("failed to init db: ", err)
	}
	sqliteDb := SQliteDB{db: database}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("failed to open channel: ", err)
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
		log.Fatal("failed to declare queue: ", err)
	}

	if err := ch.ExchangeDeclare("dlx", "direct", true, false, false, false, nil); err != nil {
		log.Fatal("failed to declare DLX exchange: ", err)
	}

	if _, err := ch.QueueDeclare("orders-dlq", true, false, false, false, nil); err != nil {
		log.Fatal("failed to declare DLQ: ", err)
	}

	if err := ch.QueueBind("orders-dlq", "", "dlx", false, nil); err != nil {
		log.Fatal("failed to bind DLQ: ", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("failed to start consumer: ", err)
	}

	log.Println("waiting for messages")

	for msg := range msgs {
		var e OrderCreatedEvent
		if err := json.Unmarshal(msg.Body, &e); err != nil {
			log.Println("failed to unmarshal message: ", err)
			msg.Nack(false, false) // send to DLQ, don't requeue infinitely
			continue
		}

		ok, err := sqliteDb.processEvent(e.EventID)
		if err != nil {
			log.Println("db error ", err)
			msg.Nack(false, true)
			continue
		}
		if !ok {
			log.Println("duplicate message ignored")
			msg.Ack(false)
			continue
		}

		log.Println("📦 order received")
		log.Println("user:", e.UserEmail)
		log.Println("product:", e.ProductID)
		log.Println("quantity:", e.Quantity)

		log.Println("Notification Sent")

		msg.Ack(false)
	}

}
