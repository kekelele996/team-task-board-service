package repository

import (
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB opens a MySQL connection and migrates all models.
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.AutoMigrate(
		&model.User{},
		&model.Workspace{},
		&model.WorkspaceMember{},
		&model.Board{},
		&model.BoardColumn{},
		&model.Task{},
		&model.Tag{},
		&model.Subtask{},
		&model.Comment{},
		&model.Attachment{},
		&model.ActivityLog{},
		&model.Notification{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return db, nil
}
