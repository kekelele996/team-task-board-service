package repository

import (
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// NotificationRepo defines notification persistence operations.
type NotificationRepo interface {
	Create(notification *model.Notification) error
	ListByUser(userID uint, limit int) ([]model.Notification, error)
	MarkRead(userID, id uint) error
	MarkAllRead(userID uint) error
	UnreadCount(userID uint) (int64, error)
}

type notificationRepo struct {
	db *gorm.DB
}

// NewNotificationRepo constructs a NotificationRepo.
func NewNotificationRepo(db *gorm.DB) NotificationRepo {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(notification *model.Notification) error {
	if err := r.db.Create(notification).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *notificationRepo) ListByUser(userID uint, limit int) ([]model.Notification, error) {
	var rows []model.Notification
	if err := r.db.Where("user_id = ?", userID).Preload("Actor").Order("id DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	return rows, nil
}

func (r *notificationRepo) MarkRead(userID, id uint) error {
	res := r.db.Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true)
	if res.Error != nil {
		return fmt.Errorf("mark notification read: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("mark notification read: %w", ErrNotFound)
	}
	return nil
}

func (r *notificationRepo) MarkAllRead(userID uint) error {
	if err := r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error; err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}

func (r *notificationRepo) UnreadCount(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}
	return count, nil
}
