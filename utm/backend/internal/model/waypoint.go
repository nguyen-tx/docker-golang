package model

import (
	"time"

	"github.com/google/uuid"
)

type Waypoint struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	X         float64   `db:"x"          json:"x"` // longitude hoặc tọa độ ngang
	Y         float64   `db:"y"          json:"y"` // latitude hoặc tọa độ dọc
	Z         float64   `db:"z"          json:"z"` // altitude (mét)
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
