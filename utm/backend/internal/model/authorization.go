package model

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizationStatus string

const (
	AuthStatusPending   AuthorizationStatus = "pending"   // Chờ xét duyệt
	AuthStatusApproved  AuthorizationStatus = "approved"  // Đã phê duyệt
	AuthStatusRejected  AuthorizationStatus = "rejected"  // Từ chối
	AuthStatusCancelled AuthorizationStatus = "cancelled" // Hủy bởi người nộp
	AuthStatusExpired   AuthorizationStatus = "expired"   // Hết hạn
	AuthStatusRevoked   AuthorizationStatus = "revoked"   // Thu hồi
)

type FlightPurpose string

const (
	PurposeCommercial   FlightPurpose = "commercial"   // Thương mại
	PurposeRecreational FlightPurpose = "recreational" // Giải trí
	PurposeSurvey       FlightPurpose = "survey"       // Khảo sát / Chụp ảnh
	PurposeDelivery     FlightPurpose = "delivery"     // Vận chuyển hàng
	PurposeEmergency    FlightPurpose = "emergency"    // Khẩn cấp
	PurposeResearch     FlightPurpose = "research"     // Nghiên cứu
)

// AuthorizationRequest - đơn xin cấp phép bay
type AuthorizationRequest struct {
	ID             uuid.UUID           `db:"id" json:"id"`
	RequestNo      string              `db:"request_no" json:"request_no"` // Số phiếu: UTM-2024-001
	PilotID        uuid.UUID           `db:"pilot_id" json:"pilot_id"`
	AircraftID     uuid.UUID           `db:"aircraft_id" json:"aircraft_id"`
	ZoneID         *uuid.UUID          `db:"zone_id" json:"zone_id,omitempty"`
	Purpose        FlightPurpose       `db:"purpose" json:"purpose"`
	Status         AuthorizationStatus `db:"status" json:"status"`
	FlightAreaJSON string              `db:"flight_area_json" json:"flight_area_json"` // GeoJSON polygon
	PlannedAltM    int                 `db:"planned_alt_m" json:"planned_alt_m"`
	PlannedStartAt time.Time           `db:"planned_start_at" json:"planned_start_at"`
	PlannedEndAt   time.Time           `db:"planned_end_at" json:"planned_end_at"`
	PilotNotes     string              `db:"pilot_notes" json:"pilot_notes,omitempty"`
	ReviewerID     *uuid.UUID          `db:"reviewer_id" json:"reviewer_id,omitempty"`
	ReviewedAt     *time.Time          `db:"reviewed_at" json:"reviewed_at,omitempty"`
	ReviewNotes    string              `db:"review_notes" json:"review_notes,omitempty"`
	ValidFrom      *time.Time          `db:"valid_from" json:"valid_from,omitempty"`
	ValidUntil     *time.Time          `db:"valid_until" json:"valid_until,omitempty"`
	CreatedAt      time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time           `db:"updated_at" json:"updated_at"`

	// Joined
	Pilot    *UserProfile `db:"-" json:"pilot,omitempty"`
	Aircraft *Aircraft    `db:"-" json:"aircraft,omitempty"`
}
