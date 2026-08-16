package repository

import (
	"errors"
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// SubtaskRepo defines subtask persistence operations.
type SubtaskRepo interface {
	Create(subtask *model.Subtask) error
	FindByID(id uint) (*model.Subtask, error)
	Update(subtask *model.Subtask) error
	Delete(id uint) error
}

type subtaskRepo struct {
	db *gorm.DB
}

// NewSubtaskRepo constructs a SubtaskRepo.
func NewSubtaskRepo(db *gorm.DB) SubtaskRepo {
	return &subtaskRepo{db: db}
}

func (r *subtaskRepo) Create(subtask *model.Subtask) error {
	if err := r.db.Create(subtask).Error; err != nil {
		return fmt.Errorf("create subtask: %w", err)
	}
	return nil
}

func (r *subtaskRepo) FindByID(id uint) (*model.Subtask, error) {
	var subtask model.Subtask
	if err := r.db.First(&subtask, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find subtask %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find subtask %d: %w", id, err)
	}
	return &subtask, nil
}

func (r *subtaskRepo) Update(subtask *model.Subtask) error {
	if err := r.db.Save(subtask).Error; err != nil {
		return fmt.Errorf("update subtask: %w", err)
	}
	return nil
}

func (r *subtaskRepo) Delete(id uint) error {
	res := r.db.Delete(&model.Subtask{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete subtask: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete subtask %d: %w", id, ErrNotFound)
	}
	return nil
}
