package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type LaneRepository struct {
	db *pgxpool.Pool
}

func NewLaneRepository(db *pgxpool.Pool) *LaneRepository {
	return &LaneRepository{db: db}
}

func scanLane(row interface{ Scan(...any) error }) (*model.Lane, error) {
	l := &model.Lane{}
	var segmentsJSON []byte
	err := row.Scan(&l.ID, &l.Name, &l.WaypointIDs, &segmentsJSON, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(segmentsJSON, &l.Segments); err != nil {
		return nil, err
	}
	return l, nil
}

func (r *LaneRepository) List(ctx context.Context) ([]*model.Lane, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, waypoint_ids, segments, created_at, updated_at FROM lanes ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Lane
	for rows.Next() {
		l, err := scanLane(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (r *LaneRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Lane, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, waypoint_ids, segments, created_at, updated_at FROM lanes WHERE id = $1`, id,
	)
	return scanLane(row)
}

func (r *LaneRepository) Create(ctx context.Context, l *model.Lane) (*model.Lane, error) {
	l.ID = uuid.New()
	segJSON, err := json.Marshal(l.Segments)
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(ctx,
		`INSERT INTO lanes (id, name, waypoint_ids, segments) VALUES ($1,$2,$3,$4)`,
		l.ID, l.Name, l.WaypointIDs, segJSON,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, l.ID)
}

func (r *LaneRepository) Update(ctx context.Context, l *model.Lane) (*model.Lane, error) {
	segJSON, err := json.Marshal(l.Segments)
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(ctx,
		`UPDATE lanes SET name=$1, waypoint_ids=$2, segments=$3, updated_at=NOW() WHERE id=$4`,
		l.Name, l.WaypointIDs, segJSON, l.ID,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, l.ID)
}

func (r *LaneRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM lanes WHERE id = $1`, id)
	return err
}
