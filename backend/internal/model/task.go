package model

import "time"

// Task is a single card on a board.
type Task struct {
	Base
	BoardID     uint       `gorm:"index;not null" json:"board_id"`
	ColumnID    uint       `gorm:"index;not null" json:"column_id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	AssigneeID  *uint      `gorm:"index" json:"assignee_id"`
	Priority    string     `gorm:"size:16;not null;default:medium" json:"priority"`
	DueDate     *time.Time `gorm:"index" json:"due_date"`
	Position    int        `gorm:"not null;default:0" json:"position"`
	CreatedBy   uint       `gorm:"not null" json:"created_by"`

	Assignee    *User        `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Column      *BoardColumn `gorm:"foreignKey:ColumnID" json:"column,omitempty"`
	Subtasks    []Subtask    `gorm:"foreignKey:TaskID" json:"subtasks,omitempty"`
	Tags        []Tag        `gorm:"many2many:task_tags;" json:"tags,omitempty"`
	Comments    []Comment    `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:TaskID" json:"attachments,omitempty"`
}
