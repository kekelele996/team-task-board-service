package model

// Comment is a discussion entry on a task.
type Comment struct {
	Base
	TaskID  uint   `gorm:"index;not null" json:"task_id"`
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	Content string `gorm:"type:text;not null" json:"content"`
	User    *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
