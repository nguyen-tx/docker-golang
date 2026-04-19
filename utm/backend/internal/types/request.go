package types

import "github.com/utm/backend/internal/model"

// ── Auth ──────────────────────────────────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	FullName     string `json:"full_name" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Organization string `json:"organization"`
	LicenseNo    string `json:"license_no"`
}

// ── Authorization ─────────────────────────────────────────────────────────────

type CreateAuthorizationRequest struct {
	AircraftID     string              `json:"aircraft_id" binding:"required,uuid"`
	ZoneID         string              `json:"zone_id"`
	Purpose        model.FlightPurpose `json:"purpose" binding:"required"`
	FlightAreaJSON string              `json:"flight_area_json" binding:"required"`
	PlannedAltM    int                 `json:"planned_alt_m" binding:"required,min=0"`
	PlannedStartAt string              `json:"planned_start_at" binding:"required"`
	PlannedEndAt   string              `json:"planned_end_at" binding:"required"`
	PilotNotes     string              `json:"pilot_notes"`
}

type ReviewAuthorizationRequest struct {
	Action      string `json:"action" binding:"required,oneof=approve reject"`
	ReviewNotes string `json:"review_notes"`
}

// ── User ──────────────────────────────────────────────────────────────────────

type UpdateUserStatusRequest struct {
	Status model.UserStatus `json:"status" binding:"required"`
}
