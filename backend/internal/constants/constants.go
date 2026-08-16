package constants

// Roles
const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// Task priorities
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// Activity action types
const (
	ActionTaskCreated   = "task_created"
	ActionTaskMoved     = "task_moved"
	ActionTaskUpdated   = "task_updated"
	ActionTaskDeleted   = "task_deleted"
	ActionMemberAdded   = "member_added"
	ActionMemberUpdated = "member_updated"
	ActionMemberRemoved = "member_removed"
	ActionCommentAdded  = "comment_added"
	ActionBoardCreated  = "board_created"
	ActionColumnCreated = "column_created"
)

// Notification types
const (
	NotifyAssigned      = "assigned"
	NotifyCommented     = "commented"
	NotifyDueSoon       = "due_soon"
	NotifyMemberInvited = "member_invited"
)

// Default board columns
var DefaultColumns = []string{"To Do", "In Progress", "Done"}

// Pagination defaults
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)
