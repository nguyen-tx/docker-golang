package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/repository/postgres"
)

type AuthFilter struct {
	PilotID string
	Status  string
}

type CreateAuthInput struct {
	AircraftID     string
	ZoneID         string
	Purpose        model.FlightPurpose
	FlightAreaJSON string
	PlannedAltM    int
	PlannedStartAt string
	PlannedEndAt   string
	PilotNotes     string
}

type AuthorizationService struct {
	repo *postgres.AuthorizationRepository
}

func NewAuthorizationService(repo *postgres.AuthorizationRepository) *AuthorizationService {
	return &AuthorizationService{repo: repo}
}

func (s *AuthorizationService) List(ctx context.Context, filter AuthFilter) ([]*model.AuthorizationRequest, error) {
	return s.repo.List(ctx, filter.PilotID, filter.Status)
}

func (s *AuthorizationService) GetByID(ctx context.Context, id string) (*model.AuthorizationRequest, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, uid)
}

func (s *AuthorizationService) Create(ctx context.Context, pilotID string, input CreateAuthInput) (*model.AuthorizationRequest, error) {
	pid, err := uuid.Parse(pilotID)
	if err != nil {
		return nil, err
	}
	aid, err := uuid.Parse(input.AircraftID)
	if err != nil {
		return nil, errors.New("invalid aircraft_id")
	}

	startAt, err := time.Parse(time.RFC3339, input.PlannedStartAt)
	if err != nil {
		return nil, errors.New("invalid planned_start_at format, use RFC3339")
	}
	endAt, err := time.Parse(time.RFC3339, input.PlannedEndAt)
	if err != nil {
		return nil, errors.New("invalid planned_end_at format, use RFC3339")
	}
	if endAt.Before(startAt) {
		return nil, errors.New("planned_end_at must be after planned_start_at")
	}

	// Đếm số đơn trong ngày để sinh request_no
	count, _ := s.repo.CountToday(ctx)
	requestNo := fmt.Sprintf("UTM-%s-%04d", time.Now().Format("20060102"), count+1)

	req := &model.AuthorizationRequest{
		ID:             uuid.New(),
		RequestNo:      requestNo,
		PilotID:        pid,
		AircraftID:     aid,
		Purpose:        input.Purpose,
		Status:         model.AuthStatusPending,
		FlightAreaJSON: input.FlightAreaJSON,
		PlannedAltM:    input.PlannedAltM,
		PlannedStartAt: startAt,
		PlannedEndAt:   endAt,
		PilotNotes:     input.PilotNotes,
	}

	if input.ZoneID != "" {
		zid, err := uuid.Parse(input.ZoneID)
		if err == nil {
			req.ZoneID = &zid
		}
	}

	return s.repo.Create(ctx, req)
}

func (s *AuthorizationService) Review(ctx context.Context, id, reviewerID, action, notes string) (*model.AuthorizationRequest, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	rid, err := uuid.Parse(reviewerID)
	if err != nil {
		return nil, err
	}

	req, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, errors.New("authorization request not found")
	}
	if req.Status != model.AuthStatusPending {
		return nil, errors.New("only pending requests can be reviewed")
	}

	now := time.Now()
	req.ReviewerID = &rid
	req.ReviewedAt = &now
	req.ReviewNotes = notes

	switch action {
	case "approve":
		req.Status = model.AuthStatusApproved
		req.ValidFrom = &req.PlannedStartAt
		req.ValidUntil = &req.PlannedEndAt
	case "reject":
		req.Status = model.AuthStatusRejected
	default:
		return nil, errors.New("invalid action")
	}

	return s.repo.Update(ctx, req)
}

func (s *AuthorizationService) Cancel(ctx context.Context, id, userID string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	req, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return errors.New("authorization request not found")
	}
	if req.PilotID.String() != userID {
		return errors.New("you can only cancel your own requests")
	}
	if req.Status != model.AuthStatusPending {
		return errors.New("only pending requests can be cancelled")
	}
	req.Status = model.AuthStatusCancelled
	_, err = s.repo.Update(ctx, req)
	return err
}
