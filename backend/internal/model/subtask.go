package model

// Subtask is a checkable item inside a task.
type Subtask struct {
	Base
	TaskID    uint   `gorm:"index;not null" json:"task_id"`
	Title     string `gorm:"size:255;not null" json:"title"`
	Completed bool   `gorm:"not null;default:false" json:"completed"`
	Position  int    `gorm:"not null;default:0" json:"position"`
}
