package model

import (
	"time"

	"github.com/google/uuid"
)

// LaneSegment mô tả hành lang bay giữa 2 waypoint liên tiếp
type LaneSegment struct {
	FromWaypointID uuid.UUID `json:"from_waypoint_id"`
	ToWaypointID   uuid.UUID `json:"to_waypoint_id"`
	Length         float64   `json:"length"`  // mét
	VMin           float64   `json:"v_min"`   // tốc độ tối thiểu (m/s)
	VMax           float64   `json:"v_max"`   // tốc độ tối đa (m/s)
}

type Lane struct {
	ID          uuid.UUID     `db:"id"           json:"id"`
	Name        string        `db:"name"         json:"name"`
	WaypointIDs []uuid.UUID   `db:"waypoint_ids" json:"waypoint_ids"` // thứ tự các waypoint
	Segments    []LaneSegment `db:"segments"     json:"segments"`     // lưu dạng JSONB
	CreatedAt   time.Time     `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"   json:"updated_at"`

	// Joined — load khi cần build SCRP request
	Waypoints []*Waypoint `db:"-" json:"waypoints,omitempty"`
}
