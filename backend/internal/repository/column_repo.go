package repository

import (
	"errors"
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// ColumnRepo defines board column persistence operations.
type ColumnRepo interface {
	Create(column *model.BoardColumn) error
	FindByID(id uint) (*model.BoardColumn, error)
	ListByBoard(boardID uint) ([]model.BoardColumn, error)
	Update(column *model.BoardColumn) error
	Delete(id uint) error
	MaxPosition(boardID uint) (int, error)
}

type columnRepo struct {
	db *gorm.DB
}

// NewColumnRepo constructs a ColumnRepo.
func NewColumnRepo(db *gorm.DB) ColumnRepo {
	return &columnRepo{db: db}
}

func (r *columnRepo) Create(column *model.BoardColumn) error {
	if err := r.db.Create(column).Error; err != nil {
		return fmt.Errorf("create column: %w", err)
	}
	return nil
}

func (r *columnRepo) FindByID(id uint) (*model.BoardColumn, error) {
	var column model.BoardColumn
	if err := r.db.First(&column, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find column %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find column %d: %w", id, err)
	}
	return &column, nil
}

func (r *columnRepo) ListByBoard(boardID uint) ([]model.BoardColumn, error) {
	var columns []model.BoardColumn
	if err := r.db.Where("board_id = ?", boardID).Order("position DESC").Find(&columns).Error; err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}
	return columns, nil
}

func (r *columnRepo) Update(column *model.BoardColumn) error {
	if err := r.db.Save(column).Error; err != nil {
		return fmt.Errorf("update column: %w", err)
	}
	return nil
}

func (r *columnRepo) Delete(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, table := range []string{"task_tags", "subtasks", "comments", "attachments"} {
			if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE task_id IN (SELECT id FROM tasks WHERE column_id = ?)", table), id).Error; err != nil {
				return fmt.Errorf("delete column %s: %w", table, err)
			}
		}
		if err := tx.Where("column_id = ?", id).Delete(&model.Task{}).Error; err != nil {
			return fmt.Errorf("delete column tasks: %w", err)
		}
		res := tx.Delete(&model.BoardColumn{}, id)
		if res.Error != nil {
			return fmt.Errorf("delete column: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("delete column %d: %w", id, ErrNotFound)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) MaxPosition(boardID uint) (int, error) {
	var max int
	err := r.db.Model(&model.BoardColumn{}).Where("board_id = ?", boardID).Select("COALESCE(MAX(position), -1)").Scan(&max).Error
	if err != nil {
		return 0, fmt.Errorf("max column position: %w", err)
	}
	return max, nil
}
