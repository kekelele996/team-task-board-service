package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles registration, login and token issuing.
type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.TokenResponse, error)
	Login(req dto.LoginRequest) (*dto.TokenResponse, error)
	ParseToken(tokenString string) (uint, error)
}

type authService struct {
	users    repository.UserRepo
	secret   string
	ttlHours int
	logger   *slog.Logger
}

// NewAuthService constructs an AuthService.
func NewAuthService(users repository.UserRepo, secret string, ttlHours int, logger *slog.Logger) AuthService {
	return &authService{users: users, secret: secret, ttlHours: ttlHours, logger: logger}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.TokenResponse, error) {
	if _, err := s.users.FindByEmail(req.Email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("register: check email: %w", err)
	}
	if _, err := s.users.FindByUsername(req.Username); err == nil {
		return nil, ErrUsernameTaken
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("register: check username: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := s.users.Create(user); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("register: create user: %w", err)
	}

	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{Token: token, User: userToDTO(user)}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.users.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("login: find user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{Token: token, User: userToDTO(user)}, nil
}

func (s *authService) ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("parse token: invalid claims")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("parse token: missing subject")
	}
	return uint(sub), nil
}

func (s *authService) issueToken(user *model.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(s.ttlHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}
	return signed, nil
}

func userToDTO(user *model.User) dto.UserDTO {
	return dto.UserDTO{ID: user.ID, Username: user.Username, Email: user.Email}
}
