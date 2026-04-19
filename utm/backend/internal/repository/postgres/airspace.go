package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type AirspaceRepository struct {
	db *pgxpool.Pool
}

func NewAirspaceRepository(db *pgxpool.Pool) *AirspaceRepository {
	return &AirspaceRepository{db: db}
}

const airspaceCols = `id, name, description, type, min_altitude_m, max_altitude_m,
	boundary_geojson, active_from, active_until, is_active, created_by, created_at, updated_at`

func scanZone(row interface {
	Scan(...any) error
}) (*model.AirspaceZone, error) {
	z := &model.AirspaceZone{}
	err := row.Scan(
		&z.ID, &z.Name, &z.Description, &z.Type,
		&z.MinAltitudeM, &z.MaxAltitudeM, &z.BoundaryGeoJSON,
		&z.ActiveFrom, &z.ActiveUntil, &z.IsActive,
		&z.CreatedBy, &z.CreatedAt, &z.UpdatedAt,
	)
	return z, err
}

func (r *AirspaceRepository) List(ctx context.Context, onlyActive bool) ([]*model.AirspaceZone, error) {
	query := `SELECT ` + airspaceCols + ` FROM airspace_zones`
	if onlyActive {
		query += ` WHERE is_active = TRUE`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.AirspaceZone
	for rows.Next() {
		z, err := scanZone(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, z)
	}
	return list, nil
}

func (r *AirspaceRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.AirspaceZone, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+airspaceCols+` FROM airspace_zones WHERE id = $1`, id)
	return scanZone(row)
}

func (r *AirspaceRepository) Create(ctx context.Context, z *model.AirspaceZone) (*model.AirspaceZone, error) {
	_, err := r.db.Exec(ctx, `
		INSERT INTO airspace_zones
		(id, name, description, type, min_altitude_m, max_altitude_m,
		 boundary_geojson, active_from, active_until, is_active, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, z.ID, z.Name, z.Description, z.Type, z.MinAltitudeM, z.MaxAltitudeM,
		z.BoundaryGeoJSON, z.ActiveFrom, z.ActiveUntil, z.IsActive, z.CreatedBy)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, z.ID)
}

func (r *AirspaceRepository) Update(ctx context.Context, z *model.AirspaceZone) (*model.AirspaceZone, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE airspace_zones SET
		name=$1, description=$2, type=$3, min_altitude_m=$4, max_altitude_m=$5,
		boundary_geojson=$6, active_from=$7, active_until=$8, is_active=$9, updated_at=NOW()
		WHERE id=$10
	`, z.Name, z.Description, z.Type, z.MinAltitudeM, z.MaxAltitudeM,
		z.BoundaryGeoJSON, z.ActiveFrom, z.ActiveUntil, z.IsActive, z.ID)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, z.ID)
}

func (r *AirspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM airspace_zones WHERE id = $1`, id)
	return err
}
