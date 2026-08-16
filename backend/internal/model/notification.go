package model

// Notification is a user-scoped in-app notice.
type Notification struct {
	Base
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	ActorID *uint  `gorm:"index" json:"actor_id"`
	Type    string `gorm:"size:32;not null" json:"type"`
	Content string `gorm:"type:text;not null" json:"content"`
	Read    bool   `gorm:"column:is_read;not null;default:true" json:"read"`
	Actor   *User  `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}
