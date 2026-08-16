package model

// WorkspaceMember links a user to a workspace with a role.
type WorkspaceMember struct {
	Base
	WorkspaceID uint       `gorm:"uniqueIndex:uk_workspace_user;not null" json:"workspace_id"`
	UserID      uint       `gorm:"uniqueIndex:uk_workspace_user;not null" json:"user_id"`
	Role        string     `gorm:"size:16;not null;default:viewer" json:"role"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Workspace   *Workspace `gorm:"foreignKey:WorkspaceID" json:"-"`
}
