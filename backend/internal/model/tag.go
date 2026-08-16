package model

// Tag is a colored label scoped to a workspace.
type Tag struct {
	Base
	WorkspaceID uint   `gorm:"index;not null" json:"workspace_id"`
	Name        string `gorm:"size:64;not null" json:"name"`
	Color       string `gorm:"size:32;not null" json:"color"`
}
