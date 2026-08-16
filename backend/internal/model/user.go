package model

// User is a registered account that can join workspaces.
type User struct {
	Base
	Username     string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email        string `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
}
