package model

import (
	"time"

	"github.com/google/uuid"
)

type ZoneType string

const (
	ZoneTypeRestricted   ZoneType = "restricted"   // Vùng cấm bay
	ZoneTypeControlled   ZoneType = "controlled"   // Vùng kiểm soát
	ZoneTypeUncontrolled ZoneType = "uncontrolled" // Vùng tự do
	ZoneTypeTFR          ZoneType = "tfr"          // Tạm thời hạn chế bay (Temporary Flight Restriction)
	ZoneTypeNFZ          ZoneType = "nfz"          // No Fly Zone
)

// Coordinate - tọa độ điểm
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// AirspaceZone - vùng không phận
type AirspaceZone struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`
	Type         ZoneType  `db:"type" json:"type"`
	MinAltitudeM int       `db:"min_altitude_m" json:"min_altitude_m"`
	MaxAltitudeM int       `db:"max_altitude_m" json:"max_altitude_m"`
	// Polygon lưu dạng GeoJSON text trong postgres
	BoundaryGeoJSON string     `db:"boundary_geojson" json:"boundary_geojson"`
	ActiveFrom      *time.Time `db:"active_from" json:"active_from,omitempty"`
	ActiveUntil     *time.Time `db:"active_until" json:"active_until,omitempty"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	CreatedBy       uuid.UUID  `db:"created_by" json:"created_by"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}
