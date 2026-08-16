package repository

import (
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// ActivityRepo defines activity log persistence operations.
type ActivityRepo interface {
	Create(activity *model.ActivityLog) error
	ListByWorkspace(workspaceID uint, limit int) ([]model.ActivityLog, error)
	ListByBoard(boardID uint, limit int) ([]model.ActivityLog, error)
}

type activityRepo struct {
	db *gorm.DB
}

// NewActivityRepo constructs an ActivityRepo.
func NewActivityRepo(db *gorm.DB) ActivityRepo {
	return &activityRepo{db: db}
}

func (r *activityRepo) Create(activity *model.ActivityLog) error {
	if err := r.db.Create(activity).Error; err != nil {
		return fmt.Errorf("create activity: %w", err)
	}
	return nil
}

func (r *activityRepo) ListByWorkspace(workspaceID uint, limit int) ([]model.ActivityLog, error) {
	var rows []model.ActivityLog
	if err := r.db.Where("workspace_id = ?", workspaceID).Preload("User").Order("id DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list workspace activities: %w", err)
	}
	return rows, nil
}

func (r *activityRepo) ListByBoard(boardID uint, limit int) ([]model.ActivityLog, error) {
	var rows []model.ActivityLog
	if err := r.db.Where("board_id = ?", boardID).Preload("User").Order("id DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list board activities: %w", err)
	}
	return rows, nil
}
