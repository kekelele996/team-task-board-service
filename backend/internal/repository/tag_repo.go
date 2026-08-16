package repository

import (
	"errors"
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// TagRepo defines tag persistence operations.
type TagRepo interface {
	Create(tag *model.Tag) error
	FindByID(id uint) (*model.Tag, error)
	ListByWorkspace(workspaceID uint) ([]model.Tag, error)
	Delete(id uint) error
}

type tagRepo struct {
	db *gorm.DB
}

// NewTagRepo constructs a TagRepo.
func NewTagRepo(db *gorm.DB) TagRepo {
	return &tagRepo{db: db}
}

func (r *tagRepo) Create(tag *model.Tag) error {
	if err := r.db.Create(tag).Error; err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	return nil
}

func (r *tagRepo) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find tag %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find tag %d: %w", id, err)
	}
	return &tag, nil
}

func (r *tagRepo) ListByWorkspace(workspaceID uint) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.Where("workspace_id = ?", workspaceID).Order("id DESC").Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	return tags, nil
}

func (r *tagRepo) Delete(id uint) error {
	res := r.db.Delete(&model.Tag{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete tag: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete tag %d: %w", id, ErrNotFound)
	}
	return nil
}
