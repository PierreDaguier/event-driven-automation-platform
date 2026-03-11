package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	dlq     string
}

type DLQEnvelope struct {
	Reason   string          `json:"reason"`
	FailedAt time.Time       `json:"failed_at"`
	Job      models.QueueJob `json:"job"`
}

func NewRabbitMQ(url, queueName, dlqName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	if _, err := ch.QueueDeclare(dlqName, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{conn: conn, channel: ch, queue: queueName, dlq: dlqName}, nil
}

func (r *RabbitMQ) Close() error {
	_ = r.channel.Close()
	return r.conn.Close()
}

func (r *RabbitMQ) Publish(ctx context.Context, job models.QueueJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return r.channel.PublishWithContext(ctx, "", r.queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
		Timestamp:    time.Now().UTC(),
	})
}

func (r *RabbitMQ) PublishDLQ(ctx context.Context, job models.QueueJob, reason string) error {
	body, err := json.Marshal(DLQEnvelope{
		Reason:   reason,
		FailedAt: time.Now().UTC(),
		Job:      job,
	})
	if err != nil {
		return err
	}
	return r.channel.PublishWithContext(ctx, "", r.dlq, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
		Timestamp:    time.Now().UTC(),
	})
}

func (r *RabbitMQ) Consume(ctx context.Context, consumerName string, handler func(context.Context, models.QueueJob) error) error {
	msgs, err := r.channel.Consume(r.queue, consumerName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			var job models.QueueJob
			if err := json.Unmarshal(msg.Body, &job); err != nil {
				_ = msg.Nack(false, false)
				continue
			}
			if err := handler(ctx, job); err != nil {
				_ = msg.Nack(false, false)
				continue
			}
			_ = msg.Ack(false)
		}
	}
}
