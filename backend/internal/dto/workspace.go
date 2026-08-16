package dto

// CreateWorkspaceRequest creates a workspace.
type CreateWorkspaceRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=128"`
	Description string `json:"description" validate:"max=2000"`
}

// UpdateWorkspaceRequest updates workspace metadata.
type UpdateWorkspaceRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=128"`
	Description string `json:"description" validate:"max=2000"`
}

// AddMemberRequest invites a member by email or username.
type AddMemberRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=3,max=64"`
	Role     string `json:"role" validate:"required,oneof=admin editor viewer"`
}

// UpdateMemberRequest changes a member role.
type UpdateMemberRequest struct {
	Role string `json:"role" validate:"required,oneof=admin editor viewer"`
}
