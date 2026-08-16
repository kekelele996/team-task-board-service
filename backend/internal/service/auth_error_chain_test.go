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

func newAuthFixture(t *testing.T) (*gorm.DB, repository.UserRepo, repository.WorkspaceRepo, repository.TaskRepo, AuthService) {
	t.Helper()
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
	userRepo := repository.NewUserRepo(db)
	wsRepo := repository.NewWorkspaceRepo(db)
	taskRepo := repository.NewTaskRepo(db)
	auth := NewAuthService(userRepo, "secret", 24, slog.Default())
	return db, userRepo, wsRepo, taskRepo, auth
}

func TestAuthErrorChain(t *testing.T) {
	_, userRepo, wsRepo, taskRepo, auth := newAuthFixture(t)

	if _, err := auth.Register(dto.RegisterRequest{Username: "alice", Email: "alice@example.com", Password: "secret123"}); err != nil {
		t.Fatalf("Register(new user) error = %v, want nil", err)
	}

	if _, err := auth.Login(dto.LoginRequest{Email: "missing@example.com", Password: "secret123"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login(unknown) error = %v, want ErrInvalidCredentials", err)
	}

	if _, err := userRepo.FindByEmail("missing@example.com"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("FindByEmail(missing) error = %v, want ErrNotFound", err)
	}
	if _, err := wsRepo.GetMember(99999, 99999); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("GetMember(missing) error = %v, want ErrNotFound", err)
	}
	if _, err := taskRepo.FindByID(99999); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("FindByID(missing) error = %v, want ErrNotFound", err)
	}
}
