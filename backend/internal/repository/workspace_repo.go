package repository

import (
	"errors"
	"fmt"

	"github.com/gbkanban/gbkanban/internal/model"
	"gorm.io/gorm"
)

// MemberDetail is a member row joined with its user profile.
type MemberDetail struct {
	model.WorkspaceMember
	Username string `json:"username"`
	Email    string `json:"email"`
}

// WorkspaceRepo defines workspace and membership persistence operations.
type WorkspaceRepo interface {
	Create(workspace *model.Workspace) error
	FindByID(id uint) (*model.Workspace, error)
	ListByUser(userID uint) ([]model.Workspace, error)
	Update(workspace *model.Workspace) error
	Delete(id uint) error

	AddMember(member *model.WorkspaceMember) error
	GetMember(workspaceID, userID uint) (*model.WorkspaceMember, error)
	ListMembers(workspaceID uint) ([]MemberDetail, error)
	UpdateMemberRole(workspaceID, userID uint, role string) error
	RemoveMember(workspaceID, userID uint) error
}

type workspaceRepo struct {
	db *gorm.DB
}

// NewWorkspaceRepo constructs a WorkspaceRepo.
func NewWorkspaceRepo(db *gorm.DB) WorkspaceRepo {
	return &workspaceRepo{db: db}
}

func (r *workspaceRepo) Create(workspace *model.Workspace) error {
	if err := r.db.Create(workspace).Error; err != nil {
		return fmt.Errorf("create workspace: %w", err)
	}
	return nil
}

func (r *workspaceRepo) FindByID(id uint) (*model.Workspace, error) {
	var workspace model.Workspace
	if err := r.db.Preload("Owner").First(&workspace, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find workspace %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find workspace %d: %w", id, err)
	}
	return &workspace, nil
}

func (r *workspaceRepo) ListByUser(userID uint) ([]model.Workspace, error) {
	var workspaces []model.Workspace
	err := r.db.Model(&model.Workspace{}).
		Joins("JOIN workspace_members ON workspace_members.workspace_id = workspaces.id").
		Where("workspace_members.user_id = ?", userID).
		Preload("Owner").
		Order("workspaces.id DESC").
		Find(&workspaces).Error
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return workspaces, nil
}

func (r *workspaceRepo) Update(workspace *model.Workspace) error {
	if err := r.db.Save(workspace).Error; err != nil {
		return fmt.Errorf("update workspace: %w", err)
	}
	return nil
}

func (r *workspaceRepo) Delete(id uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, table := range []string{"task_tags", "subtasks", "comments", "attachments"} {
			if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE task_id IN (SELECT id FROM tasks WHERE board_id IN (SELECT id FROM boards WHERE workspace_id = ?))", table), id).Error; err != nil {
				return fmt.Errorf("delete workspace %s: %w", table, err)
			}
		}
		if err := tx.Exec("DELETE FROM tasks WHERE board_id IN (SELECT id FROM boards WHERE workspace_id = ?)", id).Error; err != nil {
			return fmt.Errorf("delete workspace tasks: %w", err)
		}
		if err := tx.Exec("DELETE FROM board_columns WHERE board_id IN (SELECT id FROM boards WHERE workspace_id = ?)", id).Error; err != nil {
			return fmt.Errorf("delete workspace columns: %w", err)
		}
		if err := tx.Where("workspace_id = ?", id).Delete(&model.Board{}).Error; err != nil {
			return fmt.Errorf("delete workspace boards: %w", err)
		}
		if err := tx.Where("workspace_id = ?", id).Delete(&model.WorkspaceMember{}).Error; err != nil {
			return fmt.Errorf("delete workspace members: %w", err)
		}
		if err := tx.Where("workspace_id = ?", id).Delete(&model.ActivityLog{}).Error; err != nil {
			return fmt.Errorf("delete workspace activities: %w", err)
		}
		res := tx.Delete(&model.Workspace{}, id)
		if res.Error != nil {
			return fmt.Errorf("delete workspace: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("delete workspace %d: %w", id, ErrNotFound)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *workspaceRepo) AddMember(member *model.WorkspaceMember) error {
	if err := r.db.Create(member).Error; err != nil {
		if isDuplicate(err) {
			return fmt.Errorf("add member: %w", ErrConflict)
		}
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r *workspaceRepo) GetMember(workspaceID, userID uint) (*model.WorkspaceMember, error) {
	var member model.WorkspaceMember
	err := r.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get member: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get member: %w", err)
	}
	return &member, nil
}

func (r *workspaceRepo) ListMembers(workspaceID uint) ([]MemberDetail, error) {
	var members []MemberDetail
	err := r.db.Model(&model.WorkspaceMember{}).
		Select("workspace_members.*, users.username, users.email").
		Joins("JOIN users ON users.id = workspace_members.user_id").
		Where("workspace_members.workspace_id = ?", workspaceID).
		Order("workspace_members.id ASC").
		Scan(&members).Error
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return members, nil
}

func (r *workspaceRepo) UpdateMemberRole(workspaceID, userID uint, role string) error {
	res := r.db.Model(&model.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Update("role", role)
	if res.Error != nil {
		return fmt.Errorf("update member role: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("update member role: %w", ErrNotFound)
	}
	return nil
}

func (r *workspaceRepo) RemoveMember(workspaceID, userID uint) error {
	res := r.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).Delete(&model.WorkspaceMember{})
	if res.Error != nil {
		return fmt.Errorf("remove member: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("remove member: %w", ErrNotFound)
	}
	return nil
}
