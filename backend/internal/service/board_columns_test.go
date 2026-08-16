package service

import (
	"log/slog"
	"testing"

	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestBoardDefaultColumnsOrder(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	models := []any{
		&model.User{}, &model.Workspace{}, &model.WorkspaceMember{}, &model.Board{},
		&model.BoardColumn{}, &model.Task{}, &model.Tag{}, &model.Subtask{}, &model.Comment{},
		&model.Attachment{}, &model.ActivityLog{}, &model.Notification{},
	}
	if err := db.AutoMigrate(models...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	log := slog.Default()
	userRepo := repository.NewUserRepo(db)
	wsRepo := repository.NewWorkspaceRepo(db)
	boardRepo := repository.NewBoardRepo(db)
	columnRepo := repository.NewColumnRepo(db)
	wsService := NewWorkspaceService(wsRepo, userRepo, log)
	boardService := NewBoardService(boardRepo, columnRepo, wsRepo, wsService, log)

	owner := &model.User{Username: "alice", Email: "alice@example.com", PasswordHash: "hash"}
	if err := userRepo.Create(owner); err != nil {
		t.Fatalf("create user: %v", err)
	}
	ws, err := wsService.Create(owner.ID, dto.CreateWorkspaceRequest{Name: "w"})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	board, err := boardService.Create(owner.ID, ws.ID, dto.CreateBoardRequest{Name: "b"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}
	cols, err := boardService.ListColumns(board.ID)
	if err != nil {
		t.Fatalf("list columns: %v", err)
	}
	got := make([]string, 0, len(cols))
	for _, c := range cols {
		got = append(got, c.Name)
	}
	want := []string{"To Do", "In Progress", "Done"}
	if len(got) != len(want) {
		t.Fatalf("columns = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("columns = %v, want %v", got, want)
		}
	}
}
