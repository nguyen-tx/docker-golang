package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type VertiportRepository struct {
	db *pgxpool.Pool
}

func NewVertiportRepository(db *pgxpool.Pool) *VertiportRepository {
	return &VertiportRepository{db: db}
}

func (r *VertiportRepository) List(ctx context.Context) ([]*model.Vertiport, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, slot_duration_sec, created_at, updated_at FROM vertiports ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Vertiport
	for rows.Next() {
		v := &model.Vertiport{}
		if err := rows.Scan(&v.ID, &v.Name, &v.SlotDurationSec, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

func (r *VertiportRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Vertiport, error) {
	v := &model.Vertiport{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slot_duration_sec, created_at, updated_at FROM vertiports WHERE id = $1`, id,
	).Scan(&v.ID, &v.Name, &v.SlotDurationSec, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// load pads
	pads, err := r.listPads(ctx, id)
	if err != nil {
		return nil, err
	}
	v.Pads = pads
	return v, nil
}

func (r *VertiportRepository) Create(ctx context.Context, v *model.Vertiport) (*model.Vertiport, error) {
	v.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO vertiports (id, name, slot_duration_sec) VALUES ($1,$2,$3)`,
		v.ID, v.Name, v.SlotDurationSec,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, v.ID)
}

func (r *VertiportRepository) Update(ctx context.Context, v *model.Vertiport) (*model.Vertiport, error) {
	_, err := r.db.Exec(ctx,
		`UPDATE vertiports SET name=$1, slot_duration_sec=$2, updated_at=NOW() WHERE id=$3`,
		v.Name, v.SlotDurationSec, v.ID,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, v.ID)
}

func (r *VertiportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM vertiports WHERE id = $1`, id)
	return err
}

// --- Pads ---

func (r *VertiportRepository) listPads(ctx context.Context, vertiportID uuid.UUID) ([]*model.Pad, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, vertiport_id, name, compatible_types, created_at, updated_at FROM pads WHERE vertiport_id = $1`,
		vertiportID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Pad
	for rows.Next() {
		p := &model.Pad{}
		if err := rows.Scan(&p.ID, &p.VertiportID, &p.Name, &p.CompatibleTypes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *VertiportRepository) CreatePad(ctx context.Context, p *model.Pad) (*model.Pad, error) {
	p.ID = uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO pads (id, vertiport_id, name, compatible_types) VALUES ($1,$2,$3,$4)`,
		p.ID, p.VertiportID, p.Name, p.CompatibleTypes,
	)
	if err != nil {
		return nil, err
	}
	out := &model.Pad{}
	err = r.db.QueryRow(ctx,
		`SELECT id, vertiport_id, name, compatible_types, created_at, updated_at FROM pads WHERE id = $1`, p.ID,
	).Scan(&out.ID, &out.VertiportID, &out.Name, &out.CompatibleTypes, &out.CreatedAt, &out.UpdatedAt)
	return out, err
}

func (r *VertiportRepository) DeletePad(ctx context.Context, padID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pads WHERE id = $1`, padID)
	return err
}
