package repository

import (
	"errors"
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// BoardRepo defines board persistence operations.
type BoardRepo interface {
	Create(board *model.Board) error
	FindByID(id uint) (*model.Board, error)
	ListByWorkspace(workspaceID uint) ([]model.Board, error)
	Update(board *model.Board) error
	Delete(id uint) error
}

type boardRepo struct {
	db *gorm.DB
}

// NewBoardRepo constructs a BoardRepo.
func NewBoardRepo(db *gorm.DB) BoardRepo {
	return &boardRepo{db: db}
}

func (r *boardRepo) Create(board *model.Board) error {
	if err := r.db.Create(board).Error; err != nil {
		return fmt.Errorf("create board: %w", err)
	}
	return nil
}

func (r *boardRepo) FindByID(id uint) (*model.Board, error) {
	var board model.Board
	if err := r.db.Preload("Columns").First(&board, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find board %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find board %d: %w", id, err)
	}
	return &board, nil
}

func (r *boardRepo) ListByWorkspace(workspaceID uint) ([]model.Board, error) {
	var boards []model.Board
	if err := r.db.Where("workspace_id = ?", workspaceID).Order("id DESC").Find(&boards).Error; err != nil {
		return nil, fmt.Errorf("list boards: %w", err)
	}
	return boards, nil
}

func (r *boardRepo) Update(board *model.Board) error {
	if err := r.db.Save(board).Error; err != nil {
		return fmt.Errorf("update board: %w", err)
	}
	return nil
}

func (r *boardRepo) Delete(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, table := range []string{"task_tags", "subtasks", "comments", "attachments"} {
			if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE task_id IN (SELECT id FROM tasks WHERE board_id = ?)", table), id).Error; err != nil {
				return fmt.Errorf("delete board %s: %w", table, err)
			}
		}
		if err := tx.Where("board_id = ?", id).Delete(&model.Task{}).Error; err != nil {
			return fmt.Errorf("delete board tasks: %w", err)
		}
		if err := tx.Where("board_id = ?", id).Delete(&model.BoardColumn{}).Error; err != nil {
			return fmt.Errorf("delete board columns: %w", err)
		}
		res := tx.Delete(&model.Board{}, id)
		if res.Error != nil {
			return fmt.Errorf("delete board: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("delete board %d: %w", id, ErrNotFound)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
