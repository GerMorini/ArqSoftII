package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"search/internal/services"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rabbitmq/amqp091-go"
)

const (
	encodingJSON = "application/json"
	encodingUTF8 = "UTF-8"
)

type RabbitMQClient struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	queue      *amqp091.Queue
}

func NewRabbitMQClient(user, password, queueName, host, port string) *RabbitMQClient {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	var connection *amqp091.Connection
	var err error

	// Retry connection up to 10 times with exponential backoff
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		connection, err = amqp091.Dial(connStr)
		if err == nil {
			break
		}

		waitTime := time.Duration(i+1) * time.Second
		log.Warnf("Failed to connect to RabbitMQ (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, waitTime)
		time.Sleep(waitTime)
	}

	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ after %d attempts: %v", maxRetries, err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	// Declare queue with same settings as activities-api (durable: true)
	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable - survives broker restart
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	log.Infof("Successfully connected to RabbitMQ at %s:%s", host, port)
	return &RabbitMQClient{connection: connection, channel: channel, queue: &queue}
}

func (r *RabbitMQClient) Consume(ctx context.Context, handler func(context.Context, services.ActivityEvent) error) error {
	// Configurar el consumer
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("🎯 Consumer registered for queue: %s", r.queue.Name)

	// Loop infinito para consumir mensajes
	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Consumer context cancelled")
			return ctx.Err()

		case msg := <-msgs:
			// Deserializar mensaje
			var event services.ActivityEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Error unmarshalling message: %v", err)
				continue
			}

			// Procesar mensaje
			if err := handler(ctx, event); err != nil {
				log.Printf("❌ Error handling message: %v", err)
			}
		}
	}
}
