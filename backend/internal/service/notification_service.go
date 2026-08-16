package service

import (
	"fmt"
	"log/slog"

	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
)

// NotificationService exposes user-scoped notifications.
type NotificationService interface {
	List(userID uint, limit int) ([]model.Notification, error)
	UnreadCount(userID uint) (int64, error)
	MarkRead(userID, id uint) error
	MarkAllRead(userID uint) error
}

type notificationService struct {
	repo   repository.NotificationRepo
	logger *slog.Logger
}

// NewNotificationService constructs a NotificationService.
func NewNotificationService(repo repository.NotificationRepo, logger *slog.Logger) NotificationService {
	return &notificationService{repo: repo, logger: logger}
}

func (s *notificationService) List(userID uint, limit int) ([]model.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.ListByUser(userID, limit)
}

func (s *notificationService) UnreadCount(userID uint) (int64, error) {
	return s.repo.UnreadCount(userID)
}

func (s *notificationService) MarkRead(userID, id uint) error {
	if err := s.repo.MarkRead(userID, id); err != nil {
		return fmt.Errorf("mark notification read: %w", err)
	}
	return nil
}

func (s *notificationService) MarkAllRead(userID uint) error {
	if err := s.repo.MarkAllRead(userID); err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}
