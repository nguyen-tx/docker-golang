package model

import (
	"time"

	"github.com/google/uuid"
)

// Pad — điểm hạ cánh trong một Vertiport
type Pad struct {
	ID              uuid.UUID `db:"id"               json:"id"`
	VertiportID     uuid.UUID `db:"vertiport_id"     json:"vertiport_id"`
	Name            string    `db:"name"             json:"name"`
	CompatibleTypes []string  `db:"compatible_types" json:"compatible_types"` // loại UAV có thể hạ cánh
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"       json:"updated_at"`
}

// Vertiport — điểm hạ/cất cánh
type Vertiport struct {
	ID              uuid.UUID `db:"id"               json:"id"`
	Name            string    `db:"name"             json:"name"`
	SlotDurationSec float64   `db:"slot_duration_sec" json:"slot_duration_sec"` // giây/slot
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"       json:"updated_at"`

	// Joined
	Pads []*Pad `db:"-" json:"pads,omitempty"`
}
