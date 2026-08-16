package model

// ActivityLog records an operation inside a workspace/board.
type ActivityLog struct {
	Base
	WorkspaceID uint   `gorm:"index;not null" json:"workspace_id"`
	BoardID     uint   `gorm:"index;not null" json:"board_id"`
	TaskID      *uint  `gorm:"index" json:"task_id"`
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	Action      string `gorm:"size:64;not null" json:"action"`
	Detail      string `gorm:"type:text" json:"detail"`
	User        *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
