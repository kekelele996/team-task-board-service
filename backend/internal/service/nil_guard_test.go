package service

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func expectNoPanicErrNotFound(t *testing.T, name string, fn func() error) error {
	t.Helper()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("%s panicked: %v", name, r)
			}
		}()
		err = fn()
	}()
	return err
}

func TestMissingBoardAndTaskReturnNotFound(t *testing.T) {
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
	taskRepo := repository.NewTaskRepo(db)
	tagRepo := repository.NewTagRepo(db)
	subtaskRepo := repository.NewSubtaskRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	attachmentRepo := repository.NewAttachmentRepo(db)
	activityRepo := repository.NewActivityRepo(db)
	notifRepo := repository.NewNotificationRepo(db)
	wsService := NewWorkspaceService(wsRepo, userRepo, log)
	boardService := NewBoardService(boardRepo, columnRepo, wsRepo, wsService, log)
	taskService := NewTaskService(taskRepo, boardRepo, columnRepo, tagRepo, subtaskRepo, commentRepo, attachmentRepo, activityRepo, notifRepo, wsService, log)

	owner := &model.User{Username: "alice", Email: "alice@example.com", PasswordHash: "hash"}
	if err := userRepo.Create(owner); err != nil {
		t.Fatalf("create user: %v", err)
	}
	ws, err := wsService.Create(owner.ID, dto.CreateWorkspaceRequest{Name: "w"})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	err = expectNoPanicErrNotFound(t, "boardService.Get", func() error {
		_, err := boardService.Get(owner.ID, ws.ID, 99999)
		return err
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("boardService.Get(missing) error = %v, want ErrNotFound", err)
	}

	err = expectNoPanicErrNotFound(t, "taskService.Get", func() error {
		_, err := taskService.Get(owner.ID, ws.ID, 99999)
		return err
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("taskService.Get(missing) error = %v, want ErrNotFound", err)
	}
}
