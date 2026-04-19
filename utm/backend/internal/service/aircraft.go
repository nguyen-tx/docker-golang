package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/repository/postgres"
)

type AircraftService struct {
	repo *postgres.AircraftRepository
}

func NewAircraftService(repo *postgres.AircraftRepository) *AircraftService {
	return &AircraftService{repo: repo}
}

func (s *AircraftService) List(ctx context.Context, ownerID string) ([]*model.Aircraft, error) {
	return s.repo.List(ctx, ownerID)
}

func (s *AircraftService) GetByID(ctx context.Context, id string) (*model.Aircraft, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, uid)
}

func (s *AircraftService) Create(ctx context.Context, ownerID string, input *model.Aircraft) (*model.Aircraft, error) {
	oid, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, err
	}
	input.ID = uuid.New()
	input.OwnerID = oid
	input.Status = model.AircraftStatusActive
	return s.repo.Create(ctx, input)
}

func (s *AircraftService) Update(ctx context.Context, id string, input *model.Aircraft) (*model.Aircraft, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	input.ID = uid
	return s.repo.Update(ctx, input)
}

func (s *AircraftService) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, uid)
}
