package service

import (
	"log/slog"
	"testing"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCommentNotificationChain(t *testing.T) {
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

	alice := &model.User{Username: "alice", Email: "alice@example.com", PasswordHash: "hash"}
	bob := &model.User{Username: "bob", Email: "bob@example.com", PasswordHash: "hash"}
	if err := userRepo.Create(alice); err != nil {
		t.Fatalf("create alice: %v", err)
	}
	if err := userRepo.Create(bob); err != nil {
		t.Fatalf("create bob: %v", err)
	}
	ws, err := wsService.Create(alice.ID, dto.CreateWorkspaceRequest{Name: "w"})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	if err := wsService.AddMember(alice.ID, ws.ID, dto.AddMemberRequest{Username: "bob", Role: constants.RoleViewer}); err != nil {
		t.Fatalf("add bob: %v", err)
	}
	board, err := boardService.Create(alice.ID, ws.ID, dto.CreateBoardRequest{Name: "b"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}
	cols, _ := columnRepo.ListByBoard(board.ID)
	var todoID uint
	for _, c := range cols {
		if c.Name == "To Do" {
			todoID = c.ID
		}
	}
	task, err := taskService.Create(alice.ID, ws.ID, board.ID, dto.CreateTaskRequest{Title: "t", ColumnID: todoID, AssigneeID: &alice.ID, Priority: constants.PriorityMedium})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if _, err := taskService.CreateComment(bob.ID, ws.ID, task.ID, dto.CreateCommentRequest{Content: "hi"}); err != nil {
		t.Fatalf("create comment: %v", err)
	}

	notes, err := notifRepo.ListByUser(alice.ID, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("alice notifications = %d, want 1", len(notes))
	}
	if notes[0].ActorID == nil || *notes[0].ActorID != bob.ID {
		t.Fatalf("notification actor = %v, want bob(%d)", notes[0].ActorID, bob.ID)
	}
	unread, err := notifRepo.UnreadCount(alice.ID)
	if err != nil {
		t.Fatalf("unread count: %v", err)
	}
	if unread != 1 {
		t.Fatalf("alice unread = %d, want 1", unread)
	}
}
