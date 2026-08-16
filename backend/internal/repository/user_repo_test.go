package repository

import (
	"errors"
	"testing"

	"github.com/gbkanban/gbkanban/internal/model"
)

func TestUserRepo_CreateAndFind(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepo(db)

	tests := []struct {
		name string
		user *model.User
	}{
		{name: "alice", user: &model.User{Username: "alice", Email: "alice@example.com", PasswordHash: "hash"}},
		{name: "bob", user: &model.User{Username: "bob", Email: "bob@example.com", PasswordHash: "hash"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.Create(tt.user); err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			found, err := repo.FindByEmail(tt.user.Email)
			if err != nil {
				t.Fatalf("FindByEmail() error = %v", err)
			}
			if found.Username != tt.user.Username {
				t.Fatalf("FindByEmail() username = %s, want %s", found.Username, tt.user.Username)
			}
		})
	}
}

func TestUserRepo_DuplicateConflict(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepo(db)
	first := &model.User{Username: "alice", Email: "alice@example.com", PasswordHash: "hash"}
	if err := repo.Create(first); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	err := repo.Create(&model.User{Username: "alice2", Email: "alice@example.com", PasswordHash: "hash"})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("Create() error = %v, want ErrConflict", err)
	}
}
