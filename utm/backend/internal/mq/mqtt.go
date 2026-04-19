package mq

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/utm/backend/internal/config"
)

type MQTTPublisher struct {
	client mqtt.Client
}

func NewMQTTPublisher(cfg config.MQTTConfig) (*MQTTPublisher, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetConnectTimeout(5 * time.Second).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(_ mqtt.Client) {
			log.Info().Str("broker", cfg.Broker).Msg("mqtt connected")
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			log.Warn().Err(err).Msg("mqtt connection lost")
		})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt connect: %w", token.Error())
	}
	return &MQTTPublisher{client: client}, nil
}

// Publish gửi payload JSON lên topic với QoS 1.
func (p *MQTTPublisher) Publish(topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mqtt marshal: %w", err)
	}
	token := p.client.Publish(topic, 1, false, data)
	token.Wait()
	return token.Error()
}

func (p *MQTTPublisher) Close() {
	p.client.Disconnect(500)
}

// MQTTClient là alias của MQTTPublisher, dùng khi cần cả pub và sub sau này.
type MQTTClient = MQTTPublisher

func NewMQTTClient(cfg config.MQTTConfig) (*MQTTClient, error) {
	return NewMQTTPublisher(cfg)
}
