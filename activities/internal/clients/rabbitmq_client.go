package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMQClient struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

type ActivityEvent struct {
	Action     string `json:"action"`      // "create", "update", "delete"
	ActivityID string `json:"activity_id"`
}

func NewRabbitMQClient(host, port, user, pass, queueName string) (*RabbitMQClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	var conn *amqp.Connection
	var err error

	// Retry connection up to 10 times with exponential backoff
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		waitTime := time.Duration(i+1) * time.Second
		log.Warnf("Failed to connect to RabbitMQ (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, waitTime)
		time.Sleep(waitTime)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare queue (idempotent)
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Infof("RabbitMQ client connected to %s:%s, queue: %s", host, port, queueName)

	return &RabbitMQClient{
		conn:      conn,
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (r *RabbitMQClient) Publish(ctx context.Context, action string, activityID string) error {
	event := ActivityEvent{
		Action:     action,
		ActivityID: activityID,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = r.channel.PublishWithContext(
		publishCtx,
		"",           // exchange
		r.queueName,  // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Debugf("Published event: action=%s, activity_id=%s", action, activityID)
	return nil
}

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
