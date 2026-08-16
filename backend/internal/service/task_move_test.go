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

func newTaskMoveFixture(t *testing.T) (uint, uint, uint, uint, uint, TaskService, BoardService) {
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
	board, err := boardService.Create(owner.ID, ws.ID, dto.CreateBoardRequest{Name: "b"})
	if err != nil {
		t.Fatalf("create board: %v", err)
	}
	cols, err := columnRepo.ListByBoard(board.ID)
	if err != nil {
		t.Fatalf("list columns: %v", err)
	}
	var todoID, doneID uint
	for _, c := range cols {
		if c.Name == "To Do" {
			todoID = c.ID
		}
		if c.Name == "Done" {
			doneID = c.ID
		}
	}
	if todoID == 0 || doneID == 0 {
		t.Fatalf("default columns missing: %+v", cols)
	}
	return owner.ID, ws.ID, board.ID, todoID, doneID, taskService, boardService
}

func TestTaskMoveAndReorder(t *testing.T) {
	userID, wsID, boardID, todoID, doneID, taskService, _ := newTaskMoveFixture(t)

	mk := func(title string) *model.Task {
		t.Helper()
		task, err := taskService.Create(userID, wsID, boardID, dto.CreateTaskRequest{Title: title, ColumnID: todoID, Priority: constants.PriorityMedium})
		if err != nil {
			t.Fatalf("Create(%s) error = %v", title, err)
		}
		return task
	}
	t1 := mk("one")
	t2 := mk("two")
	t3 := mk("three")

	first, err := taskService.Get(userID, wsID, t1.ID)
	if err != nil {
		t.Fatalf("Get(first) error = %v", err)
	}
	if first.Position != 0 {
		t.Fatalf("first task position = %d, want 0", first.Position)
	}

	if _, err := taskService.Move(userID, wsID, t3.ID, dto.MoveTaskRequest{ColumnID: todoID, TaskIDs: []uint{t3.ID, t1.ID, t2.ID}}); err != nil {
		t.Fatalf("Move(reorder) error = %v", err)
	}
	list, err := taskService.List(userID, wsID, boardID, dto.ListTasksQuery{ColumnID: todoID})
	if err != nil {
		t.Fatalf("List(todo) error = %v", err)
	}
	got := []uint{list[0].ID, list[1].ID, list[2].ID}
	want := []uint{t3.ID, t1.ID, t2.ID}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("todo order = %v, want %v", got, want)
		}
	}

	if _, err := taskService.Move(userID, wsID, t1.ID, dto.MoveTaskRequest{ColumnID: doneID, Position: 0}); err != nil {
		t.Fatalf("Move(t1->done) error = %v", err)
	}
	todoList, err := taskService.List(userID, wsID, boardID, dto.ListTasksQuery{ColumnID: todoID})
	if err != nil {
		t.Fatalf("List(todo after move) error = %v", err)
	}
	for _, task := range todoList {
		if task.ID == t1.ID {
			t.Fatalf("task %d still in todo after moving to done", t1.ID)
		}
	}
	doneList, err := taskService.List(userID, wsID, boardID, dto.ListTasksQuery{ColumnID: doneID})
	if err != nil {
		t.Fatalf("List(done) error = %v", err)
	}
	found := false
	for _, task := range doneList {
		if task.ID == t1.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("task %d not found in done list", t1.ID)
	}
}
