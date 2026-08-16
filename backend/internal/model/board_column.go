package model

// BoardColumn is a single status lane on a board.
type BoardColumn struct {
	Base
	BoardID  uint   `gorm:"index;not null" json:"board_id"`
	Name     string `gorm:"size:64;not null" json:"name"`
	Position int    `gorm:"not null;default:0" json:"position"`
	Color    string `gorm:"size:32" json:"color"`
}
