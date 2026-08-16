package handler

import (
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/middleware"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
)

// WorkspaceHandler exposes workspace and member endpoints.
type WorkspaceHandler struct {
	workspaces service.WorkspaceService
}

// NewWorkspaceHandler constructs a WorkspaceHandler.
func NewWorkspaceHandler(workspaces service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaces: workspaces}
}

// List returns workspaces visible to the current user.
func (h *WorkspaceHandler) List(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	workspaces, err := h.workspaces.ListByUser(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, workspaces)
}

// Create creates a workspace and makes the creator an admin.
func (h *WorkspaceHandler) Create(c *gin.Context) {
	var req dto.CreateWorkspaceRequest
	if !bindJSON(c, &req) {
		return
	}
	workspace, err := h.workspaces.Create(middleware.CurrentUserID(c), req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, workspace)
}

// Get returns a workspace the user belongs to.
func (h *WorkspaceHandler) Get(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	workspace, err := h.workspaces.Get(middleware.CurrentUserID(c), workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, workspace)
}

// Update updates workspace metadata (admin only).
func (h *WorkspaceHandler) Update(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	var req dto.UpdateWorkspaceRequest
	if !bindJSON(c, &req) {
		return
	}
	workspace, err := h.workspaces.Update(middleware.CurrentUserID(c), workspaceID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, workspace)
}

// Delete removes a workspace (admin only).
func (h *WorkspaceHandler) Delete(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	if err := h.workspaces.Delete(middleware.CurrentUserID(c), workspaceID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// ListMembers lists workspace members.
func (h *WorkspaceHandler) ListMembers(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	members, err := h.workspaces.ListMembers(workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, members)
}

// AddMember invites a member.
func (h *WorkspaceHandler) AddMember(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	var req dto.AddMemberRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.workspaces.AddMember(middleware.CurrentUserID(c), workspaceID, req); err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, gin.H{"invited": true})
}

// UpdateMember changes a member role.
func (h *WorkspaceHandler) UpdateMember(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	targetUserID, ok := parseUintParam(c, "user_id")
	if !ok {
		return
	}
	var req dto.UpdateMemberRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.workspaces.UpdateMember(middleware.CurrentUserID(c), workspaceID, targetUserID, req); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"updated": true})
}

// RemoveMember removes a member.
func (h *WorkspaceHandler) RemoveMember(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	targetUserID, ok := parseUintParam(c, "user_id")
	if !ok {
		return
	}
	if err := h.workspaces.RemoveMember(middleware.CurrentUserID(c), workspaceID, targetUserID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"removed": true})
}
