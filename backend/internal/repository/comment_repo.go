package repository

import (
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// CommentRepo defines comment persistence operations.
type CommentRepo interface {
	Create(comment *model.Comment) error
	Delete(id uint) error
}

type commentRepo struct {
	db *gorm.DB
}

// NewCommentRepo constructs a CommentRepo.
func NewCommentRepo(db *gorm.DB) CommentRepo {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(comment *model.Comment) error {
	if err := r.db.Create(comment).Error; err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (r *commentRepo) Delete(id uint) error {
	res := r.db.Delete(&model.Comment{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete comment: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete comment %d: %w", id, ErrNotFound)
	}
	return nil
}
