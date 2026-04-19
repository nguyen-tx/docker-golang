package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type AircraftRepository struct {
	db *pgxpool.Pool
}

func NewAircraftRepository(db *pgxpool.Pool) *AircraftRepository {
	return &AircraftRepository{db: db}
}

func (r *AircraftRepository) List(ctx context.Context, ownerID string) ([]*model.Aircraft, error) {
	query := `SELECT id, owner_id, registration_no, model, manufacturer, serial_no, type, status,
		weight_kg, max_altitude_m, max_speed_kmh, max_flight_time_min, has_camera,
		insurance_no, insurance_expiry, created_at, updated_at FROM aircraft`

	args := []any{}
	if ownerID != "" {
		query += " WHERE owner_id = $1"
		args = append(args, ownerID)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Aircraft
	for rows.Next() {
		a := &model.Aircraft{}
		if err := rows.Scan(&a.ID, &a.OwnerID, &a.RegistrationNo, &a.Model, &a.Manufacturer,
			&a.SerialNo, &a.Type, &a.Status, &a.WeightKg, &a.MaxAltitudeM,
			&a.MaxSpeedKmh, &a.MaxFlightTimMin, &a.HasCamera,
			&a.InsuranceNo, &a.InsuranceExpiry, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AircraftRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Aircraft, error) {
	a := &model.Aircraft{}
	err := r.db.QueryRow(ctx, `
		SELECT id, owner_id, registration_no, model, manufacturer, serial_no, type, status,
		weight_kg, max_altitude_m, max_speed_kmh, max_flight_time_min, has_camera,
		insurance_no, insurance_expiry, created_at, updated_at
		FROM aircraft WHERE id = $1
	`, id).Scan(&a.ID, &a.OwnerID, &a.RegistrationNo, &a.Model, &a.Manufacturer,
		&a.SerialNo, &a.Type, &a.Status, &a.WeightKg, &a.MaxAltitudeM,
		&a.MaxSpeedKmh, &a.MaxFlightTimMin, &a.HasCamera,
		&a.InsuranceNo, &a.InsuranceExpiry, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *AircraftRepository) Create(ctx context.Context, a *model.Aircraft) (*model.Aircraft, error) {
	_, err := r.db.Exec(ctx, `
		INSERT INTO aircraft (id, owner_id, registration_no, model, manufacturer, serial_no,
		type, status, weight_kg, max_altitude_m, max_speed_kmh, max_flight_time_min,
		has_camera, insurance_no, insurance_expiry)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, a.ID, a.OwnerID, a.RegistrationNo, a.Model, a.Manufacturer, a.SerialNo,
		a.Type, a.Status, a.WeightKg, a.MaxAltitudeM, a.MaxSpeedKmh, a.MaxFlightTimMin,
		a.HasCamera, a.InsuranceNo, a.InsuranceExpiry)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, a.ID)
}

func (r *AircraftRepository) Update(ctx context.Context, a *model.Aircraft) (*model.Aircraft, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE aircraft SET model=$1, manufacturer=$2, status=$3, weight_kg=$4,
		max_altitude_m=$5, max_speed_kmh=$6, max_flight_time_min=$7,
		has_camera=$8, insurance_no=$9, insurance_expiry=$10, updated_at=NOW()
		WHERE id=$11
	`, a.Model, a.Manufacturer, a.Status, a.WeightKg, a.MaxAltitudeM,
		a.MaxSpeedKmh, a.MaxFlightTimMin, a.HasCamera, a.InsuranceNo, a.InsuranceExpiry, a.ID)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, a.ID)
}

func (r *AircraftRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM aircraft WHERE id = $1`, id)
	return err
}
