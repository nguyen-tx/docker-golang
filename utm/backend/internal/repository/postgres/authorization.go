package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type AuthorizationRepository struct {
	db *pgxpool.Pool
}

func NewAuthorizationRepository(db *pgxpool.Pool) *AuthorizationRepository {
	return &AuthorizationRepository{db: db}
}

func (r *AuthorizationRepository) CountToday(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM authorization_requests
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&count)
	return count, err
}

func (r *AuthorizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.AuthorizationRequest, error) {
	req := &model.AuthorizationRequest{}
	err := r.db.QueryRow(ctx, `
		SELECT id, request_no, pilot_id, aircraft_id, zone_id, purpose, status,
		flight_area_json, planned_alt_m, planned_start_at, planned_end_at, pilot_notes,
		reviewer_id, reviewed_at, review_notes, valid_from, valid_until, created_at, updated_at
		FROM authorization_requests WHERE id = $1
	`, id).Scan(
		&req.ID, &req.RequestNo, &req.PilotID, &req.AircraftID, &req.ZoneID,
		&req.Purpose, &req.Status, &req.FlightAreaJSON, &req.PlannedAltM,
		&req.PlannedStartAt, &req.PlannedEndAt, &req.PilotNotes,
		&req.ReviewerID, &req.ReviewedAt, &req.ReviewNotes,
		&req.ValidFrom, &req.ValidUntil, &req.CreatedAt, &req.UpdatedAt,
	)
	return req, err
}

func (r *AuthorizationRepository) List(ctx context.Context, pilotID, status string) ([]*model.AuthorizationRequest, error) {
	query := `SELECT id, request_no, pilot_id, aircraft_id, zone_id, purpose, status,
		flight_area_json, planned_alt_m, planned_start_at, planned_end_at, pilot_notes,
		reviewer_id, reviewed_at, review_notes, valid_from, valid_until, created_at, updated_at
		FROM authorization_requests WHERE 1=1`

	args := []any{}
	i := 1
	if pilotID != "" {
		query += " AND pilot_id = $" + string(rune('0'+i))
		args = append(args, pilotID)
		i++
	}
	if status != "" {
		query += " AND status = $" + string(rune('0'+i))
		args = append(args, status)
		i++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.AuthorizationRequest
	for rows.Next() {
		req := &model.AuthorizationRequest{}
		if err := rows.Scan(
			&req.ID, &req.RequestNo, &req.PilotID, &req.AircraftID, &req.ZoneID,
			&req.Purpose, &req.Status, &req.FlightAreaJSON, &req.PlannedAltM,
			&req.PlannedStartAt, &req.PlannedEndAt, &req.PilotNotes,
			&req.ReviewerID, &req.ReviewedAt, &req.ReviewNotes,
			&req.ValidFrom, &req.ValidUntil, &req.CreatedAt, &req.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, req)
	}
	return list, nil
}

func (r *AuthorizationRepository) Create(ctx context.Context, req *model.AuthorizationRequest) (*model.AuthorizationRequest, error) {
	_, err := r.db.Exec(ctx, `
		INSERT INTO authorization_requests
		(id, request_no, pilot_id, aircraft_id, zone_id, purpose, status,
		flight_area_json, planned_alt_m, planned_start_at, planned_end_at, pilot_notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, req.ID, req.RequestNo, req.PilotID, req.AircraftID, req.ZoneID, req.Purpose,
		req.Status, req.FlightAreaJSON, req.PlannedAltM, req.PlannedStartAt, req.PlannedEndAt, req.PilotNotes)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, req.ID)
}

func (r *AuthorizationRepository) Update(ctx context.Context, req *model.AuthorizationRequest) (*model.AuthorizationRequest, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE authorization_requests SET
		status=$1, reviewer_id=$2, reviewed_at=$3, review_notes=$4,
		valid_from=$5, valid_until=$6, updated_at=NOW()
		WHERE id=$7
	`, req.Status, req.ReviewerID, req.ReviewedAt, req.ReviewNotes,
		req.ValidFrom, req.ValidUntil, req.ID)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, req.ID)
}
