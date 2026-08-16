package service

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	users     map[string]*model.User
	createErr error
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[string]*model.User{}}
}

func (f *fakeUserRepo) Create(user *model.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	user.ID = uint(len(f.users) + 1)
	copied := *user
	f.users[user.Email] = &copied
	f.users[user.Username] = &copied
	return nil
}

func (f *fakeUserRepo) FindByID(id uint) (*model.User, error) {
	for _, user := range f.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeUserRepo) FindByEmail(email string) (*model.User, error) {
	user, ok := f.users[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (f *fakeUserRepo) FindByUsername(username string) (*model.User, error) {
	user, ok := f.users[username]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (f *fakeUserRepo) Update(user *model.User) error {
	return nil
}

func newTestAuthService(repo repository.UserRepo) AuthService {
	return NewAuthService(repo, "test-secret", 24, slog.Default())
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(repo *fakeUserRepo)
		request dto.RegisterRequest
		wantErr error
	}{
		{
			name: "success",
			request: dto.RegisterRequest{
				Username: "alice",
				Email:    "alice@example.com",
				Password: "secret123",
			},
		},
		{
			name: "duplicate email",
			setup: func(repo *fakeUserRepo) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
				repo.users["alice@example.com"] = &model.User{Base: model.Base{ID: 1}, Email: "alice@example.com", PasswordHash: string(hash)}
			},
			request: dto.RegisterRequest{
				Username: "alice2",
				Email:    "alice@example.com",
				Password: "secret123",
			},
			wantErr: ErrEmailTaken,
		},
		{
			name: "duplicate username",
			setup: func(repo *fakeUserRepo) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
				repo.users["alice"] = &model.User{Base: model.Base{ID: 1}, Username: "alice", Email: "alice@example.com", PasswordHash: string(hash)}
			},
			request: dto.RegisterRequest{
				Username: "alice",
				Email:    "alice2@example.com",
				Password: "secret123",
			},
			wantErr: ErrUsernameTaken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeUserRepo()
			if tt.setup != nil {
				tt.setup(repo)
			}
			service := newTestAuthService(repo)
			_, err := service.Register(tt.request)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Register() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	repo := newFakeUserRepo()
	repo.users["alice@example.com"] = &model.User{Base: model.Base{ID: 1}, Username: "alice", Email: "alice@example.com", PasswordHash: string(hash)}
	service := newTestAuthService(repo)

	tests := []struct {
		name    string
		request dto.LoginRequest
		wantErr error
	}{
		{
			name:    "success",
			request: dto.LoginRequest{Email: "alice@example.com", Password: "secret123"},
		},
		{
			name:    "wrong password",
			request: dto.LoginRequest{Email: "alice@example.com", Password: "wrong-pass"},
			wantErr: ErrInvalidCredentials,
		},
		{
			name:    "unknown user",
			request: dto.LoginRequest{Email: "missing@example.com", Password: "secret123"},
			wantErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(tt.request)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Login() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
