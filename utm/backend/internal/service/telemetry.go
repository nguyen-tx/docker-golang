package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/mq"
	mongoRepo "github.com/utm/backend/internal/repository/mongo"
	pgRepo "github.com/utm/backend/internal/repository/postgres"
	"github.com/utm/backend/internal/ws"
)

const mqttTopicControl = "utm/control/commands"

// ControlConfig lưu cấu hình điều khiển nhận từ FE.
type ControlConfig struct {
	mu          sync.RWMutex
	sensitivity float64
	muteUntil   time.Time
}

func newControlConfig() *ControlConfig {
	return &ControlConfig{sensitivity: 1.0}
}

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
	repo       *mongoRepo.TelemetryRepository
	authzRepo  *pgRepo.AuthorizationRepository
	cmdLogRepo *mongoRepo.CommandLogRepository
	mqtt       *mq.MQTTPublisher
	hub        *ws.Hub
	config     *ControlConfig
}

func NewTelemetryService(
	repo *mongoRepo.TelemetryRepository,
	authzRepo *pgRepo.AuthorizationRepository,
	cmdLogRepo *mongoRepo.CommandLogRepository,
	mqttPub *mq.MQTTPublisher,
) *TelemetryService {
	svc := &TelemetryService{
		repo:       repo,
		authzRepo:  authzRepo,
		cmdLogRepo: cmdLogRepo,
		mqtt:       mqttPub,
		config:     newControlConfig(),
	}
	svc.hub = ws.NewHub(svc.handleCommand)
	return svc
}

func (s *TelemetryService) Hub() *ws.Hub { return s.hub }

func (s *TelemetryService) Push(ctx context.Context, data *model.Telemetry) error {
	s.broadcastTelemetry(data)
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

func (s *TelemetryService) StreamToClient(ctx context.Context, conn *websocket.Conn) {
	s.hub.ServeClient(conn)
}

// handleCommand xử lý lệnh từ FE: cập nhật config, publish MQTT, lưu log.
func (s *TelemetryService) handleCommand(cmd ws.ControlCommand) {
	receivedAt := time.Now().UTC() // ghi nhận thời điểm nhận ngay lập tức

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
		// không thay đổi config, chỉ publish + log

	default:
		s.config.mu.Unlock()
		log.Warn().Str("action", cmd.Action).Msg("unknown control action")
		return
	}
	s.config.mu.Unlock()

	// Publish và log chạy async — không chặn readPump
	go s.publishAndLog(cmd, receivedAt)
}

// publishAndLog gửi lệnh lên MQTT broker và lưu log vào MongoDB.
func (s *TelemetryService) publishAndLog(cmd ws.ControlCommand, receivedAt time.Time) {
	entry := &model.ControlCommandLog{
		Action:     cmd.Action,
		Params:     string(cmd.Params),
		ReceivedAt: receivedAt,
		MQTTTopic:  mqttTopicControl,
	}

	if s.mqtt != nil {
		if err := s.mqtt.Publish(mqttTopicControl, cmd); err != nil {
			log.Error().Err(err).Str("action", cmd.Action).Msg("mqtt publish failed")
		} else {
			entry.PublishedAt = time.Now().UTC()
			log.Info().
				Str("action", cmd.Action).
				Str("topic", mqttTopicControl).
				Time("received_at", receivedAt).
				Time("published_at", entry.PublishedAt).
				Dur("publish_latency", entry.PublishedAt.Sub(receivedAt)).
				Msg("control command published")
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
		log.Warn().Err(err).Str("auth_request_id", t.AuthRequestID).Msg("lookup auth request failed")
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
	payload, err := json.Marshal(alert)
	if err != nil {
		log.Error().Err(err).Msg("marshal lane alert failed")
		return
	}
	s.hub.Broadcast(payload)
}

func (s *TelemetryService) broadcastTelemetry(data *model.Telemetry) {
	type telemetryMsg struct {
		Type    string           `json:"type"`
		Payload *model.Telemetry `json:"payload"`
	}
	payload, err := json.Marshal(telemetryMsg{Type: "telemetry", Payload: data})
	if err != nil {
		log.Error().Err(err).Msg("marshal telemetry failed")
		return
	}
	s.hub.Broadcast(payload)
}
