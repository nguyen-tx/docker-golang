package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type WaypointRepository struct {
	db *pgxpool.Pool
}

func NewWaypointRepository(db *pgxpool.Pool) *WaypointRepository {
	return &WaypointRepository{db: db}
}

const waypointCols = `id, name, x, y, z, created_at, updated_at`

func scanWaypoint(row interface{ Scan(...any) error }) (*model.Waypoint, error) {
	w := &model.Waypoint{}
	err := row.Scan(&w.ID, &w.Name, &w.X, &w.Y, &w.Z, &w.CreatedAt, &w.UpdatedAt)
	return w, err
}

func (r *WaypointRepository) List(ctx context.Context) ([]*model.Waypoint, error) {
	rows, err := r.db.Query(ctx, `SELECT `+waypointCols+` FROM waypoints ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Waypoint
	for rows.Next() {
		w, err := scanWaypoint(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	return list, nil
}

func (r *WaypointRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Waypoint, error) {
	row := r.db.QueryRow(ctx, `SELECT `+waypointCols+` FROM waypoints WHERE id = $1`, id)
	return scanWaypoint(row)
}

// FindByIDs trả về map[id]*Waypoint để build LaneProto
func (r *WaypointRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*model.Waypoint, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]*model.Waypoint{}, nil
	}

	// build $1,$2,... placeholders
	args := make([]any, len(ids))
	placeholders := ""
	for i, id := range ids {
		args[i] = id
		if i > 0 {
			placeholders += ","
		}
		placeholders += "$" + string(rune('1'+i))
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+waypointCols+` FROM waypoints WHERE id = ANY($1)`,
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]*model.Waypoint, len(ids))
	for rows.Next() {
		w, err := scanWaypoint(rows)
		if err != nil {
			return nil, err
		}
		result[w.ID] = w
	}
	return result, nil
}

func (r *WaypointRepository) Create(ctx context.Context, w *model.Waypoint) (*model.Waypoint, error) {
	w.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO waypoints (id, name, x, y, z) VALUES ($1,$2,$3,$4,$5)`,
		w.ID, w.Name, w.X, w.Y, w.Z,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, w.ID)
}

func (r *WaypointRepository) Update(ctx context.Context, w *model.Waypoint) (*model.Waypoint, error) {
	_, err := r.db.Exec(ctx,
		`UPDATE waypoints SET name=$1, x=$2, y=$3, z=$4, updated_at=NOW() WHERE id=$5`,
		w.Name, w.X, w.Y, w.Z, w.ID,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, w.ID)
}

func (r *WaypointRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM waypoints WHERE id = $1`, id)
	return err
}
