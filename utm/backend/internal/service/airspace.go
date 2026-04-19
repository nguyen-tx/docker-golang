package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/repository/postgres"
)

type AirspaceService struct {
	repo *postgres.AirspaceRepository
}

func NewAirspaceService(repo *postgres.AirspaceRepository) *AirspaceService {
	return &AirspaceService{repo: repo}
}

func (s *AirspaceService) List(ctx context.Context, onlyActive bool) ([]*model.AirspaceZone, error) {
	return s.repo.List(ctx, onlyActive)
}

func (s *AirspaceService) GetByID(ctx context.Context, id string) (*model.AirspaceZone, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	return s.repo.FindByID(ctx, uid)
}

func (s *AirspaceService) Create(ctx context.Context, createdBy string, input *model.AirspaceZone) (*model.AirspaceZone, error) {
	uid, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, err
	}
	input.ID = uuid.New()
	input.CreatedBy = uid
	input.IsActive = true
	return s.repo.Create(ctx, input)
}

func (s *AirspaceService) Update(ctx context.Context, id string, input *model.AirspaceZone) (*model.AirspaceZone, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	input.ID = uid
	return s.repo.Update(ctx, input)
}

func (s *AirspaceService) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.Delete(ctx, uid)
}
