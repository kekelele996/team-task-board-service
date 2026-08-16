package repository

import (
	"errors"
	"testing"

	"github.com/gbkanban/gbkanban/internal/model"
)

func TestTaskRepo_ListFilterAndMove(t *testing.T) {
	db := newTestDB(t)
	taskRepo := NewTaskRepo(db)
	user := seedUser(t, db, "alice", "alice@example.com")
	workspace := seedWorkspace(t, db, user.ID)
	board := seedBoard(t, db, workspace.ID)
	todo := seedColumn(t, db, board.ID, "To Do", 0)
	done := seedColumn(t, db, board.ID, "Done", 1)

	task := &model.Task{BoardID: board.ID, ColumnID: todo.ID, Title: "ship api", Priority: "high", CreatedBy: user.ID}
	if err := taskRepo.Create(task); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	byColumn, err := taskRepo.ListByBoard(board.ID, TaskFilter{ColumnID: todo.ID})
	if err != nil {
		t.Fatalf("ListByBoard() error = %v", err)
	}
	if len(byColumn) != 1 {
		t.Fatalf("ListByBoard() len = %d, want 1", len(byColumn))
	}

	byKeyword, err := taskRepo.ListByBoard(board.ID, TaskFilter{Keyword: "ship"})
	if err != nil {
		t.Fatalf("ListByBoard keyword error = %v", err)
	}
	if len(byKeyword) != 1 {
		t.Fatalf("ListByBoard keyword len = %d, want 1", len(byKeyword))
	}

	if err := taskRepo.UpdatePosition(task.ID, done.ID, 3); err != nil {
		t.Fatalf("UpdatePosition() error = %v", err)
	}
	moved, err := taskRepo.FindByID(task.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if moved.ColumnID != done.ID || moved.Position != 3 {
		t.Fatalf("moved task = column %d position %d, want %d/3", moved.ColumnID, moved.Position, done.ID)
	}

	if err := taskRepo.Delete(task.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err = taskRepo.FindByID(task.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("FindByID() after delete error = %v, want ErrNotFound", err)
	}
}
