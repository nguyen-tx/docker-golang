package model

import (
	"time"

	"github.com/google/uuid"
)

type AircraftType string

const (
	AircraftTypeFixedWing  AircraftType = "fixed_wing"  // Cánh cứng
	AircraftTypeMultiRotor AircraftType = "multi_rotor" // Nhiều cánh quạt
	AircraftTypeHelicopter AircraftType = "helicopter"  // Trực thăng
	AircraftTypeHybrid     AircraftType = "hybrid"      // Hybrid
)

type AircraftStatus string

const (
	AircraftStatusActive      AircraftStatus = "active"
	AircraftStatusInactive    AircraftStatus = "inactive"
	AircraftStatusMaintenance AircraftStatus = "maintenance"
	AircraftStatusGrounded    AircraftStatus = "grounded" // Bị đình chỉ
)

// Aircraft - thông tin đăng ký máy bay không người lái (UAV/Drone)
type Aircraft struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	OwnerID         uuid.UUID      `db:"owner_id" json:"owner_id"`
	RegistrationNo  string         `db:"registration_no" json:"registration_no"` // Số đăng ký
	Model           string         `db:"model" json:"model"`
	Manufacturer    string         `db:"manufacturer" json:"manufacturer"`
	SerialNo        string         `db:"serial_no" json:"serial_no"`
	Type            AircraftType   `db:"type" json:"type"`
	Status          AircraftStatus `db:"status" json:"status"`
	WeightKg        float64        `db:"weight_kg" json:"weight_kg"`
	MaxAltitudeM    int            `db:"max_altitude_m" json:"max_altitude_m"` // Độ cao tối đa (mét)
	MaxSpeedKmh     int            `db:"max_speed_kmh" json:"max_speed_kmh"`
	MaxFlightTimMin int            `db:"max_flight_time_min" json:"max_flight_time_min"` // Thời gian bay tối đa (phút)
	HasCamera       bool           `db:"has_camera" json:"has_camera"`
	InsuranceNo     string         `db:"insurance_no" json:"insurance_no,omitempty"`
	InsuranceExpiry *time.Time     `db:"insurance_expiry" json:"insurance_expiry,omitempty"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`

	// Joined
	Owner *UserProfile `db:"-" json:"owner,omitempty"`
}
