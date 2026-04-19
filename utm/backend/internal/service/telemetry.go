package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/mq"
	mongoRepo "github.com/utm/backend/internal/repository/mongo"
	pgRepo "github.com/utm/backend/internal/repository/postgres"
	"github.com/utm/backend/internal/ws"
)

const mqttTopicControl = "utm/control/commands"

type ControlConfig struct {
	mu          sync.RWMutex
	sensitivity float64
	muteUntil   time.Time
}

func newControlConfig() *ControlConfig { return &ControlConfig{sensitivity: 1.0} }

func (c *ControlConfig) isMuted() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Now().Before(c.muteUntil)
}

func (c *ControlConfig) getSensitivity() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sensitivity
}

type TelemetryService struct {
	repo        *mongoRepo.TelemetryRepository
	authzRepo   *pgRepo.AuthorizationRepository
	cmdLogRepo  *mongoRepo.CommandLogRepository
	mqtt        *mq.MQTTClient
	alertClient *ws.AlertClient // WS client kết nối đến utm-alert
	config      *ControlConfig
}

func NewTelemetryService(
	repo *mongoRepo.TelemetryRepository,
	authzRepo *pgRepo.AuthorizationRepository,
	cmdLogRepo *mongoRepo.CommandLogRepository,
	mqttClient *mq.MQTTClient,
	alertURL string, // ws://utm-alert:8081/ws/backend
) *TelemetryService {
	svc := &TelemetryService{
		repo:       repo,
		authzRepo:  authzRepo,
		cmdLogRepo: cmdLogRepo,
		mqtt:       mqttClient,
		config:     newControlConfig(),
	}
	svc.alertClient = ws.NewAlertClient(alertURL, svc.handleRawCommand)
	return svc
}

// AlertClient trả về client để router gọi Run().
func (s *TelemetryService) AlertClient() *ws.AlertClient { return s.alertClient }

func (s *TelemetryService) Push(ctx context.Context, data *model.Telemetry) error {
	go s.checkAndBroadcastAlert(context.Background(), data)
	go func() {
		if err := s.repo.Insert(context.Background(), data); err != nil {
			log.Error().Err(err).Msg("insert telemetry failed")
		}
	}()
	return nil
}

func (s *TelemetryService) GetSessionHistory(ctx context.Context, sessionID string) ([]*model.Telemetry, error) {
	return s.repo.FindBySession(ctx, sessionID)
}

func (s *TelemetryService) GetActiveSessions(ctx context.Context) ([]*model.FlightSession, error) {
	return s.repo.FindActiveSessions(ctx)
}

// handleRawCommand nhận raw JSON từ utm-alert (control command do FE gửi).
func (s *TelemetryService) handleRawCommand(raw []byte) {
	receivedAt := time.Now().UTC()

	var cmd struct {
		Type   string          `json:"type"`
		Action string          `json:"action"`
		Params json.RawMessage `json:"params,omitempty"`
	}
	if err := json.Unmarshal(raw, &cmd); err != nil || cmd.Type != "control" {
		return
	}

	s.config.mu.Lock()
	switch cmd.Action {
	case "mute_alert":
		var p struct {
			DurationMs int `json:"duration_ms"`
		}
		if err := json.Unmarshal(cmd.Params, &p); err != nil || p.DurationMs <= 0 {
			s.config.mu.Unlock()
			log.Warn().Msg("mute_alert: missing duration_ms")
			return
		}
		s.config.muteUntil = time.Now().Add(time.Duration(p.DurationMs) * time.Millisecond)

	case "set_sensitivity":
		var p struct {
			Value float64 `json:"value"`
		}
		if err := json.Unmarshal(cmd.Params, &p); err != nil || p.Value <= 0 || p.Value > 1 {
			s.config.mu.Unlock()
			log.Warn().Msg("set_sensitivity: value must be 0.0–1.0")
			return
		}
		s.config.sensitivity = p.Value

	case "acknowledge":
		// không thay đổi config

	default:
		s.config.mu.Unlock()
		log.Warn().Str("action", cmd.Action).Msg("unknown control action")
		return
	}
	s.config.mu.Unlock()

	go s.publishAndLog(cmd.Action, cmd.Params, receivedAt)
}

func (s *TelemetryService) publishAndLog(action string, params json.RawMessage, receivedAt time.Time) {
	entry := &model.ControlCommandLog{
		Action:     action,
		Params:     string(params),
		ReceivedAt: receivedAt,
		MQTTTopic:  mqttTopicControl,
	}
	if s.mqtt != nil {
		payload := map[string]any{"type": "control", "action": action, "params": params}
		if err := s.mqtt.Publish(mqttTopicControl, payload); err != nil {
			log.Error().Err(err).Msg("mqtt publish control failed")
		} else {
			entry.PublishedAt = time.Now().UTC()
			log.Info().
				Str("action", action).
				Time("received_at", receivedAt).
				Dur("publish_latency", entry.PublishedAt.Sub(receivedAt)).
				Msg("control command published to mqtt")
		}
	}
	if err := s.cmdLogRepo.Insert(context.Background(), entry); err != nil {
		log.Error().Err(err).Msg("save command log failed")
	}
}

func (s *TelemetryService) checkAndBroadcastAlert(ctx context.Context, t *model.Telemetry) {
	if s.config.isMuted() {
		return
	}
	if t.AuthRequestID == "" {
		return
	}
	id, err := uuid.Parse(t.AuthRequestID)
	if err != nil {
		return
	}
	authz, err := s.authzRepo.FindByID(ctx, id)
	if err != nil {
		log.Warn().Err(err).Msg("lookup auth request failed")
		return
	}
	if authz.FlightAreaJSON == "" {
		return
	}
	scaledAltM := int(float64(authz.PlannedAltM) / s.config.getSensitivity())
	alert := checkDeviation(t, authz.FlightAreaJSON, scaledAltM)
	if alert == nil {
		return
	}
	s.sendToAlert(alert)
}

// sendToAlert serialize và gửi lane alert đến utm-alert qua WS.
func (s *TelemetryService) sendToAlert(alert *model.LaneAlert) {
	data, err := json.Marshal(alert)
	if err != nil {
		log.Error().Err(err).Msg("marshal alert message failed")
		return
	}
	s.alertClient.Send(data)
}
