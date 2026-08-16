package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/middleware"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TaskHandler exposes task, subtask, comment, attachment, tag, activity and stats endpoints.
type TaskHandler struct {
	tasks     service.TaskService
	uploadDir string
}

// NewTaskHandler constructs a TaskHandler.
func NewTaskHandler(tasks service.TaskService, uploadDir string) *TaskHandler {
	return &TaskHandler{tasks: tasks, uploadDir: uploadDir}
}

// ListTasks lists tasks with filters and search.
func (h *TaskHandler) ListTasks(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	var query dto.ListTasksQuery
	if !bindQuery(c, &query) {
		return
	}
	tasks, err := h.tasks.List(middleware.CurrentUserID(c), workspaceID, boardID, query)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, tasks)
}

// CreateTask creates a task.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	var req dto.CreateTaskRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.Create(middleware.CurrentUserID(c), workspaceID, boardID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, task)
}

// GetTask returns full task detail.
func (h *TaskHandler) GetTask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	taskID, ok := parseUintParam(c, "task_id")
	if !ok {
		return
	}
	task, err := h.tasks.Get(middleware.CurrentUserID(c), workspaceID, taskID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, task)
}

// UpdateTask updates task fields.
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	taskID, ok := parseUintParam(c, "task_id")
	if !ok {
		return
	}
	var req dto.UpdateTaskRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.Update(middleware.CurrentUserID(c), workspaceID, taskID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, task)
}

// DeleteTask deletes a task.
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	taskID, ok := parseUintParam(c, "task_id")
	if !ok {
		return
	}
	if err := h.tasks.Delete(middleware.CurrentUserID(c), workspaceID, taskID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// MoveTask handles drag reorder and cross-column status changes.
func (h *TaskHandler) MoveTask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	taskID, ok := parseUintParam(c, "task_id")
	if !ok {
		return
	}
	var req dto.MoveTaskRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.Move(middleware.CurrentUserID(c), workspaceID, taskID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, task)
}

// ListTags lists workspace tags.
func (h *TaskHandler) ListTags(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	tags, err := h.tasks.ListTags(workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, tags)
}

// CreateTag creates a tag.
func (h *TaskHandler) CreateTag(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	var req dto.CreateTagRequest
	if !bindJSON(c, &req) {
		return
	}
	tag, err := h.tasks.CreateTag(middleware.CurrentUserID(c), workspaceID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, tag)
}

// DeleteTag deletes a tag.
func (h *TaskHandler) DeleteTag(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	tagID, ok := parseUintParam(c, "tag_id")
	if !ok {
		return
	}
	if err := h.tasks.DeleteTag(middleware.CurrentUserID(c), workspaceID, tagID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// CreateSubtask adds a subtask.
func (h *TaskHandler) CreateSubtask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	taskID, ok := parseUintParam(c, "task_id")
	if !ok {
		return
	}
	var req dto.CreateSubtaskRequest
	if !bindJSON(c, &req) {
		return
	}
	subtask, err := h.tasks.CreateSubtask(middleware.CurrentUserID(c), workspaceID, taskID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, subtask)
}

// UpdateSubtask toggles or renames a subtask.
func (h *TaskHandler) UpdateSubtask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	subtaskID, ok := parseUintParam(c, "subtask_id")
	if !ok {
		return
	}
	var req dto.UpdateSubtaskRequest
	if !bindJSON(c, &req) {
		return
	}
	subtask, err := h.tasks.UpdateSubtask(middleware.CurrentUserID(c), workspaceID, subtaskID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, subtask)
}

// DeleteSubtask deletes a subtask.
func (h *TaskHandler) DeleteSubtask(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	subtaskID, ok := parseUintParam(c, "subtask_id")
	if !ok {
		return
	}
	if err := h.tasks.DeleteSubtask(middleware.CurrentUserID(c), workspaceID, subtaskID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// CreateComment adds a comment.
func (h *TaskHandler) CreateComment(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	taskID, ok := parseUintParam(c, "task_id")
	if !ok {
		return
	}
	var req dto.CreateCommentRequest
	if !bindJSON(c, &req) {
		return
	}
	comment, err := h.tasks.CreateComment(middleware.CurrentUserID(c), workspaceID, taskID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, comment)
}

// CreateAttachment uploads a file to a task.
func (h *TaskHandler) CreateAttachment(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	taskID, ok := parseUintParam(c, "task_id")
	if !ok {
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, constants.CodeBadRequest, "missing file field")
		return
	}
	src, err := fileHeader.Open()
	if err != nil {
		respondError(c, fmt.Errorf("open upload: %w", err))
		return
	}
	defer src.Close()

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		respondError(c, fmt.Errorf("create upload dir: %w", err))
		return
	}
	ext := filepath.Ext(fileHeader.Filename)
	name := uuid.NewString() + ext
	path := filepath.Join(h.uploadDir, name)
	dst, err := os.Create(path)
	if err != nil {
		respondError(c, fmt.Errorf("create upload file: %w", err))
		return
	}
	size, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil {
		respondError(c, fmt.Errorf("write upload: %w", copyErr))
		return
	}
	if closeErr != nil {
		respondError(c, fmt.Errorf("close upload: %w", closeErr))
		return
	}

	attachment := &model.Attachment{
		FileName:    fileHeader.Filename,
		FilePath:    "/uploads/" + name,
		Size:        size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}
	saved, err := h.tasks.CreateAttachment(middleware.CurrentUserID(c), workspaceID, taskID, attachment)
	if err != nil {
		_ = os.Remove(path)
		respondError(c, err)
		return
	}
	response.Created(c, saved)
}

// DeleteAttachment deletes an attachment.
func (h *TaskHandler) DeleteAttachment(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	attachmentID, ok := parseUintParam(c, "attachment_id")
	if !ok {
		return
	}
	if err := h.tasks.DeleteAttachment(middleware.CurrentUserID(c), workspaceID, attachmentID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// BoardActivities lists board activity logs.
func (h *TaskHandler) BoardActivities(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	limit := defaultLimit(c)
	activities, err := h.tasks.ListBoardActivities(workspaceID, boardID, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, activities)
}

// WorkspaceActivities lists workspace activity logs.
func (h *TaskHandler) WorkspaceActivities(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	limit := defaultLimit(c)
	activities, err := h.tasks.ListWorkspaceActivities(workspaceID, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, activities)
}

// Stats returns the workspace statistics view.
func (h *TaskHandler) Stats(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	stats, err := h.tasks.Stats(workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, stats)
}

func defaultLimit(c *gin.Context) int {
	limit := 50
	if raw := c.Query("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return limit
}
