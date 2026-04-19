package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Telemetry - dữ liệu vị trí/trạng thái máy bay real-time (lưu MongoDB)
type Telemetry struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AircraftID     string        `bson:"aircraft_id" json:"aircraft_id"`
	SessionID      string        `bson:"session_id" json:"session_id"` // ID phiên bay
	AuthRequestID  string        `bson:"auth_request_id,omitempty" json:"auth_request_id,omitempty"`
	Latitude       float64       `bson:"lat" json:"lat"`
	Longitude      float64       `bson:"lng" json:"lng"`
	AltitudeM      float64       `bson:"altitude_m" json:"altitude_m"`
	SpeedKmh       float64       `bson:"speed_kmh" json:"speed_kmh"`
	HeadingDeg     float64       `bson:"heading_deg" json:"heading_deg"` // Hướng (0-360)
	BatteryPct     int           `bson:"battery_pct" json:"battery_pct"` // Pin (%)
	SignalStrength int           `bson:"signal_strength" json:"signal_strength"`
	IsOnGround     bool          `bson:"is_on_ground" json:"is_on_ground"`
	Timestamp      time.Time     `bson:"timestamp" json:"timestamp"`
}

// FlightSession - phiên bay (lưu MongoDB)
type FlightSession struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID     string        `bson:"session_id" json:"session_id"`
	AircraftID    string        `bson:"aircraft_id" json:"aircraft_id"`
	PilotID       string        `bson:"pilot_id" json:"pilot_id"`
	AuthRequestID string        `bson:"auth_request_id,omitempty" json:"auth_request_id,omitempty"`
	StartedAt     time.Time     `bson:"started_at" json:"started_at"`
	EndedAt       *time.Time    `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
	MaxAltitudeM  float64       `bson:"max_altitude_m" json:"max_altitude_m"`
	TotalDistKm   float64       `bson:"total_dist_km" json:"total_dist_km"`
	IsActive      bool          `bson:"is_active" json:"is_active"`
}
