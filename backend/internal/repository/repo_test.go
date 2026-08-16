package repository

import (
	"fmt"
	"testing"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
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
		t.Fatalf("migrate sqlite: %v", err)
	}
	return db
}

func seedUser(t *testing.T, db *gorm.DB, username, email string) *model.User {
	t.Helper()
	user := &model.User{Username: username, Email: email, PasswordHash: "hash"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return user
}

func seedWorkspace(t *testing.T, db *gorm.DB, ownerID uint) *model.Workspace {
	t.Helper()
	workspace := &model.Workspace{Name: "workspace", OwnerID: ownerID}
	if err := db.Create(workspace).Error; err != nil {
		t.Fatalf("seed workspace: %v", err)
	}
	return workspace
}

func seedBoard(t *testing.T, db *gorm.DB, workspaceID uint) *model.Board {
	t.Helper()
	board := &model.Board{WorkspaceID: workspaceID, Name: "board", CreatedBy: 1}
	if err := db.Create(board).Error; err != nil {
		t.Fatalf("seed board: %v", err)
	}
	return board
}

func seedColumn(t *testing.T, db *gorm.DB, boardID uint, name string, position int) *model.BoardColumn {
	t.Helper()
	column := &model.BoardColumn{BoardID: boardID, Name: name, Position: position}
	if err := db.Create(column).Error; err != nil {
		t.Fatalf("seed column: %v", err)
	}
	return column
}
