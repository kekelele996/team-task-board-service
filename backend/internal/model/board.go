package model

// Board belongs to a workspace and contains columns.
type Board struct {
	Base
	WorkspaceID uint          `gorm:"index;not null" json:"workspace_id"`
	Name        string        `gorm:"size:128;not null" json:"name"`
	Description string        `gorm:"type:text" json:"description"`
	CreatedBy   uint          `gorm:"not null" json:"created_by"`
	Columns     []BoardColumn `gorm:"foreignKey:BoardID" json:"columns,omitempty"`
}
