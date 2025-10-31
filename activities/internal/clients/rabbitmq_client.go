package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

// RabbitMQClient implementa el publisher/consumer con reintentos
type RabbitMQClient struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

type ActivityEvent struct {
	Action     string `json:"action"`
	ActivityID string `json:"activity_id"`
}

// NewRabbitMQClient intenta conectar con reintentos exponenciales
func NewRabbitMQClient(host, port, user, pass, queueName string) (*RabbitMQClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	var conn *amqp.Connection
	var err error

	maxRetries := 6
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		wait := time.Duration((i+1)*2) * time.Second
		log.Warnf("RabbitMQ connect attempt %d/%d failed: %v - retrying in %v", i+1, maxRetries, err, wait)
		time.Sleep(wait)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Infof("Connected to RabbitMQ %s:%s queue=%s", host, port, queueName)
	return &RabbitMQClient{conn: conn, channel: ch, queueName: queueName}, nil
}

// Publish publica un evento de actividad
func (r *RabbitMQClient) Publish(ctx context.Context, action, activityID string) error {
	ev := ActivityEvent{Action: action, ActivityID: activityID}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return r.channel.PublishWithContext(pubCtx, "", r.queueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         b,
		DeliveryMode: amqp.Persistent,
	})
}

// Consume arranca la escucha y entrega los eventos al handler
func (r *RabbitMQClient) Consume(ctx context.Context, handler func(ctx context.Context, action, activityID string) error) error {
	msgs, err := r.channel.Consume(r.queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var ev ActivityEvent
			if err := json.Unmarshal(d.Body, &ev); err != nil {
				log.Warnf("invalid rabbit message: %v", err)
				continue
			}
			_ = handler(ctx, ev.Action, ev.ActivityID)
		}
	}()
	return nil
}

// Close cierra canal y conexión
func (r *RabbitMQClient) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return err
		}
	}
	log.Info("RabbitMQ client closed")
	return nil
}
