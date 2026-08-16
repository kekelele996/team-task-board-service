package model

// Workspace is the top-level collaboration container.
type Workspace struct {
	Base
	Name        string `gorm:"size:128;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	OwnerID     uint   `gorm:"index;not null" json:"owner_id"`
	Owner       *User  `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}
