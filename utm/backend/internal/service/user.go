package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/repository/postgres"
)

type UserService struct {
	repo *postgres.UserRepository
}

func NewUserService(repo *postgres.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) List(ctx context.Context) ([]*model.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, uid)
}

func (s *UserService) UpdateStatus(ctx context.Context, id string, status model.UserStatus) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, uid, status)
}
