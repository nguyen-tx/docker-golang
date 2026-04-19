package model

import "time"

type AlertSeverity string

const (
	AlertSeverityNormal  AlertSeverity = "normal"
	AlertSeverityWarning AlertSeverity = "warning" // lệch nhẹ: 0–50m
	AlertSeverityDanger  AlertSeverity = "danger"  // lệch nghiêm trọng: >50m hoặc vượt altitude
)

// LaneAlert là bản tin gửi xuống FE khi drone bay lệch khỏi vùng được cấp phép.
type LaneAlert struct {
	Type        string        `json:"type"`         // luôn là "lane_alert"
	AircraftID  string        `json:"aircraft_id"`
	SessionID   string        `json:"session_id"`
	Severity    AlertSeverity `json:"severity"`
	DeviationM  float64       `json:"deviation_m"`  // khoảng cách tới biên gần nhất (m)
	Lat         float64       `json:"lat"`
	Lng         float64       `json:"lng"`
	AltitudeM   float64       `json:"altitude_m"`
	Message     string        `json:"message"`
	Timestamp   time.Time     `json:"timestamp"`
}
