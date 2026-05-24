package model

import (
	"time"

	"github.com/google/uuid"
)

type FlightNotificationStatus string

const (
	FlightStatusPending      FlightNotificationStatus = "pending"       // vừa tạo, đang chờ SCRP
	FlightStatusSoftReserved FlightNotificationStatus = "soft_reserved" // SCRP approved, chờ operator xác nhận
	FlightStatusCommitted    FlightNotificationStatus = "committed"     // operator xác nhận thực hiện
	FlightStatusUTMRejected  FlightNotificationStatus = "utm_rejected"  // SCRP từ chối
	FlightStatusDCSRejected  FlightNotificationStatus = "dcs_rejected"  // operator từ chối phương án UTM
)

type FlightNotification struct {
	ID                   uuid.UUID                `db:"id"                    json:"id"`
	OperatorID           uuid.UUID                `db:"operator_id"           json:"operator_id"`
	DroneID              string                   `db:"drone_id"              json:"drone_id"`
	UavType              string                   `db:"uav_type"              json:"uav_type"`
	LaneID               uuid.UUID                `db:"lane_id"               json:"lane_id"`
	VWaypoints           []float64                `db:"v_waypoints"           json:"v_waypoints"`           // vận tốc tại từng waypoint (m/s)
	DestinationVertiport string                   `db:"destination_vertiport" json:"destination_vertiport"` // vertiport_id hạ cánh
	TDes                 float64                  `db:"t_des"                 json:"t_des"`                 // thời điểm cất cánh mong muốn (unix)
	TTakeoff             float64                  `db:"t_takeoff"             json:"t_takeoff"`              // thời điểm cất cánh (unix)
	TLandEstimated       float64                  `db:"t_land_estimated"      json:"t_land_estimated"`       // ước tính thời điểm hạ cánh (unix)
	Priority             int32                    `db:"priority"              json:"priority"`
	Soc0                 float64                  `db:"soc_0"                 json:"soc_0"`    // % pin ban đầu
	CBat                 float64                  `db:"c_bat"                 json:"c_bat"`    // dung lượng pin (Wh)
	PHover               float64                  `db:"p_hover"               json:"p_hover"`  // công suất hover (W)
	Status               FlightNotificationStatus `db:"status"                json:"status"`
	CreatedAt            time.Time                `db:"created_at"            json:"created_at"`
	UpdatedAt            time.Time                `db:"updated_at"            json:"updated_at"`

	// Joined khi cần
	Lane *Lane `db:"-" json:"lane,omitempty"`
}

// ApprovedPlan — lưu kết quả SCRP approve, độc lập với FlightNotification
// dùng để build approved_plans cho SCRP request tiếp theo
type ApprovedPlan struct {
	ID             uuid.UUID           `db:"id"              json:"id"`
	NotificationID uuid.UUID           `db:"notification_id" json:"notification_id"`
	TDepStar       float64             `db:"t_dep_star"      json:"t_dep_star"`       // thời điểm cất cánh được giao (unix)
	TLandAssigned  float64             `db:"t_land_assigned" json:"t_land_assigned"`  // thời điểm hạ cánh được giao (unix)
	PadID          string              `db:"pad_id"          json:"pad_id"`
	SlotIndex      int32               `db:"slot_index"      json:"slot_index"`
	WaypointTimes  []WaypointTimeEntry `db:"waypoint_times"  json:"waypoint_times"`   // JSONB
	ExpiresAt      float64             `db:"expires_at"      json:"expires_at"`        // unix
	DelaySeconds   float64             `db:"delay_seconds"   json:"delay_seconds"`
	DelaySource    string              `db:"delay_source"    json:"delay_source"`
	EnergyEstimate float64             `db:"energy_estimate" json:"energy_estimate"`
	SocMeaning     float64             `db:"soc_meaning"     json:"soc_meaning"`
	CreatedAt      time.Time           `db:"created_at"      json:"created_at"`
}

type WaypointTimeEntry struct {
	WaypointID string  `json:"waypoint_id"`
	Time       float64 `json:"time"` // unix
}
