package service

import (
	"context"

	"github.com/utm/backend/internal/model"
	mongoRepo "github.com/utm/backend/internal/repository/mongo"
)

type NotificationService struct {
	repo *mongoRepo.NotificationRepository
}

func NewNotificationService(repo *mongoRepo.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) ListForUser(ctx context.Context, userID string, unreadOnly bool) ([]*model.Notification, error) {
	return s.repo.FindByRecipient(ctx, userID, unreadOnly)
}

func (s *NotificationService) MarkRead(ctx context.Context, id string) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *NotificationService) Send(ctx context.Context, n *model.Notification) error {
	return s.repo.Insert(ctx, n)
}
