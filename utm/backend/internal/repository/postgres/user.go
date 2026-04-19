package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utm/backend/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, full_name, phone, role, status, license_no, organization)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, user.ID, user.Email, user.PasswordHash, user.FullName, user.Phone,
		user.Role, user.Status, user.LicenseNo, user.Organization)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, full_name, phone, role, status, license_no, organization, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone,
		&user.Role, &user.Status, &user.LicenseNo, &user.Organization, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, full_name, phone, role, status, license_no, organization, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone,
		&user.Role, &user.Status, &user.LicenseNo, &user.Organization, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*model.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, email, full_name, phone, role, status, license_no, organization, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Phone, &u.Role,
			&u.Status, &u.LicenseNo, &u.Organization, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.UserStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}
