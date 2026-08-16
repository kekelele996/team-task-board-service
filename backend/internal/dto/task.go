package dto

import "time"

// CreateTaskRequest creates a task inside a column.
type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=20000"`
	ColumnID    uint       `json:"column_id" validate:"required"`
	AssigneeID  *uint      `json:"assignee_id"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=high medium low"`
	DueDate     *time.Time `json:"due_date"`
	TagIDs      []uint     `json:"tag_ids"`
}

// UpdateTaskRequest edits task fields.
type UpdateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=20000"`
	AssigneeID  *uint      `json:"assignee_id"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=high medium low"`
	DueDate     *time.Time `json:"due_date"`
	TagIDs      []uint     `json:"tag_ids"`
}

// MoveTaskRequest updates column and position for drag-and-drop.
type MoveTaskRequest struct {
	ColumnID uint   `json:"column_id" validate:"required"`
	Position int    `json:"position" validate:"gte=0"`
	TaskIDs  []uint `json:"task_ids" validate:"omitempty"`
}

// CreateSubtaskRequest adds a subtask.
type CreateSubtaskRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// UpdateSubtaskRequest toggles a subtask.
type UpdateSubtaskRequest struct {
	Title     string `json:"title" validate:"required,min=1,max=255"`
	Completed bool   `json:"completed"`
}

// CreateCommentRequest adds a comment.
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

// CreateTagRequest creates a workspace tag.
type CreateTagRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=64"`
	Color string `json:"color" validate:"required,max=32"`
}

// ListTasksQuery carries board filters and search.
type ListTasksQuery struct {
	Pagination
	ColumnID   uint   `form:"column_id"`
	AssigneeID *uint  `form:"assignee_id"`
	Priority   string `form:"priority"`
	TagID      uint   `form:"tag_id"`
	DueBefore  string `form:"due_before"`
	DueAfter   string `form:"due_after"`
	Keyword    string `form:"keyword"`
}
