package repository

import (
	"errors"
	"testing"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/model"
)

func TestWorkspaceRepo_Members(t *testing.T) {
	db := newTestDB(t)
	repo := NewWorkspaceRepo(db)
	owner := seedUser(t, db, "owner", "owner@example.com")
	member := seedUser(t, db, "member", "member@example.com")
	workspace := seedWorkspace(t, db, owner.ID)

	if err := repo.AddMember(&model.WorkspaceMember{WorkspaceID: workspace.ID, UserID: owner.ID, Role: constants.RoleAdmin}); err != nil {
		t.Fatalf("AddMember owner error = %v", err)
	}
	if err := repo.AddMember(&model.WorkspaceMember{WorkspaceID: workspace.ID, UserID: member.ID, Role: constants.RoleViewer}); err != nil {
		t.Fatalf("AddMember member error = %v", err)
	}

	err := repo.AddMember(&model.WorkspaceMember{WorkspaceID: workspace.ID, UserID: member.ID, Role: constants.RoleEditor})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("duplicate AddMember error = %v, want ErrConflict", err)
	}

	members, err := repo.ListMembers(workspace.ID)
	if err != nil {
		t.Fatalf("ListMembers() error = %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("ListMembers() len = %d, want 2", len(members))
	}

	if err := repo.UpdateMemberRole(workspace.ID, member.ID, constants.RoleEditor); err != nil {
		t.Fatalf("UpdateMemberRole() error = %v", err)
	}
	updated, err := repo.GetMember(workspace.ID, member.ID)
	if err != nil {
		t.Fatalf("GetMember() error = %v", err)
	}
	if updated.Role != constants.RoleEditor {
		t.Fatalf("GetMember() role = %s, want editor", updated.Role)
	}
}
