package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type HandlerFunc func(body []byte) error

type Subscriber struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewSubscriber(uri string) (*Subscriber, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	return &Subscriber{conn: conn, channel: ch}, nil
}

func (s *Subscriber) Subscribe(queue string, handler HandlerFunc) error {
	_, err := s.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	msgs, err := s.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue consume: %w", err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Error().Err(err).Str("queue", queue).Msg("message handler error")
				d.Nack(false, true)
				continue
			}
			d.Ack(false)
		}
	}()

	log.Info().Str("queue", queue).Msg("subscribed")
	return nil
}

func (s *Subscriber) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
