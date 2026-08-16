package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// TaskFilter describes optional filters for listing tasks on a board.
type TaskFilter struct {
	ColumnID   uint
	AssigneeID *uint
	Priority   string
	TagID      uint
	DueBefore  *time.Time
	DueAfter   *time.Time
	Keyword    string
}

// TaskRepo defines task persistence operations.
type TaskRepo interface {
	Create(task *model.Task) error
	FindByID(id uint) (*model.Task, error)
	FindDetailed(id uint) (*model.Task, error)
	ListByBoard(boardID uint, filter TaskFilter) ([]model.Task, error)
	Update(task *model.Task) error
	Delete(id uint) error
	MaxPosition(boardID, columnID uint) (int, error)
	UpdatePosition(taskID, columnID uint, position int) error
	CountByColumn(workspaceID uint) ([]ColumnCount, error)
	CountByAssignee(workspaceID uint) ([]AssigneeCount, error)
	CompletedInRange(workspaceID uint, from, to time.Time) ([]DailyCount, error)
	ListOverdue(workspaceID uint) ([]model.Task, error)
	SetTagged(task *model.Task, tagIDs []uint) error
}

// ColumnCount is a count grouped by column name.
type ColumnCount struct {
	ColumnID uint   `json:"column_id"`
	Name     string `json:"name"`
	Count    int64  `json:"count"`
}

// AssigneeCount is a count grouped by assignee.
type AssigneeCount struct {
	AssigneeID uint   `json:"assignee_id"`
	Username   string `json:"username"`
	Count      int64  `json:"count"`
}

// DailyCount is the number of completed tasks on a single day.
type DailyCount struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

type taskRepo struct {
	db *gorm.DB
}

// NewTaskRepo constructs a TaskRepo.
func NewTaskRepo(db *gorm.DB) TaskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(task *model.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (r *taskRepo) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find task %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find task %d: %w", id, err)
	}
	return &task, nil
}

func (r *taskRepo) FindDetailed(id uint) (*model.Task, error) {
	var task model.Task
	err := r.db.Preload("Assignee").Preload("Column").
		Preload("Subtasks").Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB { return db.Order("comments.id ASC") }).
		Preload("Comments.User").
		Preload("Attachments").
		First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find task %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find task %d: %w", id, err)
	}
	return &task, nil
}

func (r *taskRepo) ListByBoard(boardID uint, filter TaskFilter) ([]model.Task, error) {
	q := r.db.Model(&model.Task{}).Where("tasks.board_id = ?", boardID)
	if filter.ColumnID != 0 {
		q = q.Where("tasks.column_id = ?", filter.ColumnID)
	}
	if filter.AssigneeID != nil {
		q = q.Where("tasks.assignee_id = ?", *filter.AssigneeID)
	}
	if filter.Priority != "" {
		q = q.Where("tasks.priority = ?", filter.Priority)
	}
	if filter.TagID != 0 {
		q = q.Joins("JOIN task_tags ON task_tags.task_id = tasks.id AND task_tags.tag_id = ?", filter.TagID)
	}
	if filter.DueBefore != nil {
		q = q.Where("tasks.due_date IS NOT NULL AND tasks.due_date <= ?", *filter.DueBefore)
	}
	if filter.DueAfter != nil {
		q = q.Where("tasks.due_date IS NOT NULL AND tasks.due_date >= ?", *filter.DueAfter)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		q = q.Where("tasks.title LIKE ? OR tasks.description LIKE ?", like, like)
	}
	var tasks []model.Task
	if err := q.Preload("Assignee").Preload("Tags").Order("tasks.id ASC").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return tasks, nil
}

func (r *taskRepo) Update(task *model.Task) error {
	if err := r.db.Save(task).Error; err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (r *taskRepo) Delete(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, table := range []string{"task_tags", "subtasks", "comments", "attachments"} {
			if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE task_id = ?", table), id).Error; err != nil {
				return fmt.Errorf("delete task %s: %w", table, err)
			}
		}
		res := tx.Delete(&model.Task{}, id)
		if res.Error != nil {
			return fmt.Errorf("delete task: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("delete task %d: %w", id, ErrNotFound)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepo) MaxPosition(boardID, columnID uint) (int, error) {
	var max int
	err := r.db.Model(&model.Task{}).
		Where("board_id = ? AND column_id = ?", boardID, columnID).
		Select("COALESCE(MAX(position), 0)").Scan(&max).Error
	if err != nil {
		return 0, fmt.Errorf("max task position: %w", err)
	}
	return max, nil
}

func (r *taskRepo) UpdatePosition(taskID, columnID uint, position int) error {
	res := r.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]any{
		"column_id": columnID,
		"position":  position,
	})
	if res.Error != nil {
		return fmt.Errorf("update task position: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("update task position: %w", ErrNotFound)
	}
	return nil
}

func (r *taskRepo) SetTagged(task *model.Task, tagIDs []uint) error {
	if err := r.db.Model(task).Association("Tags").Replace(tagsFromIDs(tagIDs)); err != nil {
		return fmt.Errorf("replace task tags: %w", err)
	}
	return nil
}

func tagsFromIDs(ids []uint) []model.Tag {
	tags := make([]model.Tag, 0, len(ids))
	for _, id := range ids {
		tags = append(tags, model.Tag{Base: model.Base{ID: id}})
	}
	return tags
}

func (r *taskRepo) CountByColumn(workspaceID uint) ([]ColumnCount, error) {
	var rows []ColumnCount
	err := r.db.Model(&model.Task{}).
		Select("tasks.column_id, board_columns.name, COUNT(tasks.id) as count").
		Joins("JOIN boards ON boards.id = tasks.board_id").
		Joins("JOIN board_columns ON board_columns.id = tasks.column_id").
		Where("boards.workspace_id = ?", workspaceID).
		Group("tasks.column_id, board_columns.name").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("count by column: %w", err)
	}
	return rows, nil
}

func (r *taskRepo) CountByAssignee(workspaceID uint) ([]AssigneeCount, error) {
	var rows []AssigneeCount
	err := r.db.Model(&model.Task{}).
		Select("tasks.assignee_id, users.username, COUNT(tasks.id) as count").
		Joins("JOIN boards ON boards.id = tasks.board_id").
		Joins("JOIN users ON users.id = tasks.assignee_id").
		Where("boards.workspace_id = ? AND tasks.assignee_id IS NOT NULL", workspaceID).
		Group("tasks.assignee_id, users.username").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("count by assignee: %w", err)
	}
	return rows, nil
}

func (r *taskRepo) CompletedInRange(workspaceID uint, from, to time.Time) ([]DailyCount, error) {
	var rows []DailyCount
	err := r.db.Model(&model.Task{}).
		Select("DATE(tasks.updated_at) as day, COUNT(tasks.id) as count").
		Joins("JOIN boards ON boards.id = tasks.board_id").
		Joins("JOIN board_columns ON board_columns.id = tasks.column_id").
		Where("boards.workspace_id = ? AND board_columns.name = ? AND tasks.updated_at >= ? AND tasks.updated_at <= ?", workspaceID, "Done", from, to).
		Group("DATE(tasks.updated_at)").
		Order("day ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("completed in range: %w", err)
	}
	return rows, nil
}

func (r *taskRepo) ListOverdue(workspaceID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Model(&model.Task{}).
		Joins("JOIN boards ON boards.id = tasks.board_id").
		Joins("JOIN board_columns ON board_columns.id = tasks.column_id").
		Where("boards.workspace_id = ? AND tasks.due_date IS NOT NULL AND tasks.due_date < ? AND board_columns.name <> ?", workspaceID, time.Now(), "Done").
		Preload("Assignee").
		Order("tasks.due_date ASC").
		Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("list overdue tasks: %w", err)
	}
	return tasks, nil
}
