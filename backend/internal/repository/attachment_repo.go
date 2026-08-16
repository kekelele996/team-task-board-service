package repository

import (
	"errors"
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// AttachmentRepo defines attachment persistence operations.
type AttachmentRepo interface {
	Create(attachment *model.Attachment) error
	FindByID(id uint) (*model.Attachment, error)
	Delete(id uint) error
}

type attachmentRepo struct {
	db *gorm.DB
}

// NewAttachmentRepo constructs an AttachmentRepo.
func NewAttachmentRepo(db *gorm.DB) AttachmentRepo {
	return &attachmentRepo{db: db}
}

func (r *attachmentRepo) Create(attachment *model.Attachment) error {
	if err := r.db.Create(attachment).Error; err != nil {
		return fmt.Errorf("create attachment: %w", err)
	}
	return nil
}

func (r *attachmentRepo) FindByID(id uint) (*model.Attachment, error) {
	var attachment model.Attachment
	if err := r.db.First(&attachment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find attachment %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find attachment %d: %w", id, err)
	}
	return &attachment, nil
}

func (r *attachmentRepo) Delete(id uint) error {
	res := r.db.Delete(&model.Attachment{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete attachment: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete attachment %d: %w", id, ErrNotFound)
	}
	return nil
}
