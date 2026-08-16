package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
)

// WorkspaceService handles workspace and membership business logic.
type WorkspaceService interface {
	Create(userID uint, req dto.CreateWorkspaceRequest) (*model.Workspace, error)
	ListByUser(userID uint) ([]model.Workspace, error)
	Get(userID, workspaceID uint) (*model.Workspace, error)
	Update(userID, workspaceID uint, req dto.UpdateWorkspaceRequest) (*model.Workspace, error)
	Delete(userID, workspaceID uint) error

	ListMembers(workspaceID uint) ([]repository.MemberDetail, error)
	AddMember(actorID, workspaceID uint, req dto.AddMemberRequest) error
	UpdateMember(actorID, workspaceID, targetUserID uint, req dto.UpdateMemberRequest) error
	RemoveMember(actorID, workspaceID, targetUserID uint) error

	RoleOf(workspaceID, userID uint) (string, error)
	RequireRole(workspaceID, userID uint, roles ...string) error
}

type workspaceService struct {
	workspaces repository.WorkspaceRepo
	users      repository.UserRepo
	logger     *slog.Logger
}

// NewWorkspaceService constructs a WorkspaceService.
func NewWorkspaceService(workspaces repository.WorkspaceRepo, users repository.UserRepo, logger *slog.Logger) WorkspaceService {
	return &workspaceService{workspaces: workspaces, users: users, logger: logger}
}

func (s *workspaceService) Create(userID uint, req dto.CreateWorkspaceRequest) (*model.Workspace, error) {
	workspace := &model.Workspace{Name: req.Name, Description: req.Description, OwnerID: userID}
	if err := s.workspaces.Create(workspace); err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}
	if err := s.workspaces.AddMember(&model.WorkspaceMember{
		WorkspaceID: workspace.ID,
		UserID:      userID,
		Role:        constants.RoleAdmin,
	}); err != nil {
		return nil, fmt.Errorf("create workspace owner membership: %w", err)
	}
	return workspace, nil
}

func (s *workspaceService) ListByUser(userID uint) ([]model.Workspace, error) {
	return s.workspaces.ListByUser(userID)
}

func (s *workspaceService) Get(userID, workspaceID uint) (*model.Workspace, error) {
	if err := s.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer); err != nil {
		return nil, err
	}
	return s.workspaces.FindByID(workspaceID)
}

func (s *workspaceService) Update(userID, workspaceID uint, req dto.UpdateWorkspaceRequest) (*model.Workspace, error) {
	if err := s.RequireRole(workspaceID, userID, constants.RoleAdmin); err != nil {
		return nil, err
	}
	workspace, err := s.workspaces.FindByID(workspaceID)
	if err != nil {
		return nil, err
	}
	workspace.Name = req.Name
	workspace.Description = req.Description
	if err := s.workspaces.Update(workspace); err != nil {
		return nil, fmt.Errorf("update workspace: %w", err)
	}
	return workspace, nil
}

func (s *workspaceService) Delete(userID, workspaceID uint) error {
	if err := s.RequireRole(workspaceID, userID, constants.RoleAdmin); err != nil {
		return err
	}
	if err := s.workspaces.Delete(workspaceID); err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}
	return nil
}

func (s *workspaceService) ListMembers(workspaceID uint) ([]repository.MemberDetail, error) {
	return s.workspaces.ListMembers(workspaceID)
}

func (s *workspaceService) AddMember(actorID, workspaceID uint, req dto.AddMemberRequest) error {
	if err := s.RequireRole(workspaceID, actorID, constants.RoleAdmin); err != nil {
		return err
	}
	var target *model.User
	var err error
	if req.Email != "" {
		target, err = s.users.FindByEmail(req.Email)
	} else if req.Username != "" {
		target, err = s.users.FindByUsername(req.Username)
	} else {
		return ErrInvalidMember
	}
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrInvalidMember
		}
		return fmt.Errorf("add member: find user: %w", err)
	}
	member := &model.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      target.ID,
		Role:        req.Role,
	}
	if err := s.workspaces.AddMember(member); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return repository.ErrConflict
		}
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (s *workspaceService) UpdateMember(actorID, workspaceID, targetUserID uint, req dto.UpdateMemberRequest) error {
	if err := s.RequireRole(workspaceID, actorID, constants.RoleAdmin); err != nil {
		return err
	}
	if err := s.workspaces.UpdateMemberRole(workspaceID, targetUserID, req.Role); err != nil {
		return fmt.Errorf("update member: %w", err)
	}
	return nil
}

func (s *workspaceService) RemoveMember(actorID, workspaceID, targetUserID uint) error {
	if err := s.RequireRole(workspaceID, actorID, constants.RoleAdmin); err != nil {
		return err
	}
	if actorID == targetUserID {
		return ErrForbidden
	}
	if err := s.workspaces.RemoveMember(workspaceID, targetUserID); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

func (s *workspaceService) RoleOf(workspaceID, userID uint) (string, error) {
	member, err := s.workspaces.GetMember(workspaceID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrForbidden
		}
		return "", fmt.Errorf("get role: %w", err)
	}
	return member.Role, nil
}

func (s *workspaceService) RequireRole(workspaceID, userID uint, roles ...string) error {
	role, err := s.RoleOf(workspaceID, userID)
	if err != nil {
		return err
	}
	for _, allowed := range roles {
		if role == allowed {
			return nil
		}
	}
	return ErrForbidden
}
