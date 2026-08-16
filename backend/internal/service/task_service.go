package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
)

// TaskService handles task, subtask, comment and attachment business logic.
type TaskService interface {
	Create(userID, workspaceID, boardID uint, req dto.CreateTaskRequest) (*model.Task, error)
	List(userID, workspaceID, boardID uint, query dto.ListTasksQuery) ([]model.Task, error)
	Get(userID, workspaceID, taskID uint) (*model.Task, error)
	Update(userID, workspaceID, taskID uint, req dto.UpdateTaskRequest) (*model.Task, error)
	Delete(userID, workspaceID, taskID uint) error
	Move(userID, workspaceID, taskID uint, req dto.MoveTaskRequest) (*model.Task, error)

	CreateTag(userID, workspaceID uint, req dto.CreateTagRequest) (*model.Tag, error)
	ListTags(workspaceID uint) ([]model.Tag, error)
	DeleteTag(userID, workspaceID, tagID uint) error

	CreateSubtask(userID, workspaceID, taskID uint, req dto.CreateSubtaskRequest) (*model.Subtask, error)
	UpdateSubtask(userID, workspaceID, subtaskID uint, req dto.UpdateSubtaskRequest) (*model.Subtask, error)
	DeleteSubtask(userID, workspaceID, subtaskID uint) error

	CreateComment(userID, workspaceID, taskID uint, req dto.CreateCommentRequest) (*model.Comment, error)

	CreateAttachment(userID, workspaceID, taskID uint, attachment *model.Attachment) (*model.Attachment, error)
	DeleteAttachment(userID, workspaceID, attachmentID uint) error

	ListBoardActivities(workspaceID, boardID uint, limit int) ([]model.ActivityLog, error)
	ListWorkspaceActivities(workspaceID uint, limit int) ([]model.ActivityLog, error)

	Stats(workspaceID uint) (map[string]any, error)
}

type taskService struct {
	tasks         repository.TaskRepo
	boards        repository.BoardRepo
	columns       repository.ColumnRepo
	tags          repository.TagRepo
	subtasks      repository.SubtaskRepo
	comments      repository.CommentRepo
	attachments   repository.AttachmentRepo
	activities    repository.ActivityRepo
	notifications repository.NotificationRepo
	workspace     WorkspaceService
	logger        *slog.Logger
}

// NewTaskService constructs a TaskService.
func NewTaskService(
	tasks repository.TaskRepo,
	boards repository.BoardRepo,
	columns repository.ColumnRepo,
	tags repository.TagRepo,
	subtasks repository.SubtaskRepo,
	comments repository.CommentRepo,
	attachments repository.AttachmentRepo,
	activities repository.ActivityRepo,
	notifications repository.NotificationRepo,
	workspace WorkspaceService,
	logger *slog.Logger,
) TaskService {
	return &taskService{
		tasks:         tasks,
		boards:        boards,
		columns:       columns,
		tags:          tags,
		subtasks:      subtasks,
		comments:      comments,
		attachments:   attachments,
		activities:    activities,
		notifications: notifications,
		workspace:     workspace,
		logger:        logger,
	}
}

func (s *taskService) requireEditor(workspaceID, userID uint) error {
	return s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor)
}

func (s *taskService) requireViewer(workspaceID, userID uint) error {
	return s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer)
}

func (s *taskService) ensureBoard(workspaceID, boardID uint) error {
	board, err := s.boards.FindByID(boardID)
	if err != nil {
		return err
	}
	if board.WorkspaceID != workspaceID {
		return repository.ErrNotFound
	}
	return nil
}

func (s *taskService) taskWorkspaceID(taskID uint) (uint, error) {
	task, err := s.tasks.FindByID(taskID)
	if err != nil {
		return 0, err
	}
	if task == nil {
		return 0, repository.ErrNotFound
	}
	board, err := s.boards.FindByID(task.BoardID)
	if err != nil {
		return 0, err
	}
	return board.WorkspaceID, nil
}

func (s *taskService) Create(userID, workspaceID, boardID uint, req dto.CreateTaskRequest) (*model.Task, error) {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return nil, err
	}
	if err := s.ensureBoard(workspaceID, boardID); err != nil {
		return nil, err
	}
	column, err := s.columns.FindByID(req.ColumnID)
	if err != nil {
		return nil, err
	}
	if column.BoardID != boardID {
		return nil, repository.ErrNotFound
	}
	max, err := s.tasks.MaxPosition(boardID, req.ColumnID)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	priority := req.Priority
	if priority == "" {
		priority = constants.PriorityMedium
	}
	task := &model.Task{
		BoardID:     boardID,
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
		Priority:    priority,
		DueDate:     req.DueDate,
		Position:    max + 1,
		CreatedBy:   userID,
	}
	if err := s.tasks.Create(task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	if err := s.replaceTags(workspaceID, task, req.TagIDs); err != nil {
		return nil, err
	}
	s.logActivity(workspaceID, boardID, &task.ID, userID, constants.ActionTaskCreated, fmt.Sprintf("创建了任务「%s」", task.Title))
	if req.AssigneeID != nil && *req.AssigneeID != userID {
		s.notify(*req.AssigneeID, userID, constants.NotifyAssigned, fmt.Sprintf("任务「%s」已分配给你", task.Title))
	}
	return s.tasks.FindDetailed(task.ID)
}

func (s *taskService) List(userID, workspaceID, boardID uint, query dto.ListTasksQuery) ([]model.Task, error) {
	if err := s.requireViewer(workspaceID, userID); err != nil {
		return nil, err
	}
	if err := s.ensureBoard(workspaceID, boardID); err != nil {
		return nil, err
	}
	filter := repository.TaskFilter{
		ColumnID:   query.ColumnID,
		AssigneeID: query.AssigneeID,
		Priority:   query.Priority,
		TagID:      query.TagID,
		Keyword:    query.Keyword,
	}
	if query.DueBefore != "" {
		t, err := time.Parse("2006-01-02", query.DueBefore)
		if err != nil {
			return nil, fmt.Errorf("parse due_before: %w", err)
		}
		t = t.Add(24*time.Hour - time.Second)
		filter.DueBefore = &t
	}
	if query.DueAfter != "" {
		t, err := time.Parse("2006-01-02", query.DueAfter)
		if err != nil {
			return nil, fmt.Errorf("parse due_after: %w", err)
		}
		filter.DueAfter = &t
	}
	return s.tasks.ListByBoard(boardID, filter)
}

func (s *taskService) Get(userID, workspaceID, taskID uint) (*model.Task, error) {
	if err := s.requireViewer(workspaceID, userID); err != nil {
		return nil, err
	}
	wsID, err := s.taskWorkspaceID(taskID)
	if err != nil {
		return nil, err
	}
	if wsID != workspaceID {
		return nil, repository.ErrNotFound
	}
	return s.tasks.FindDetailed(taskID)
}

func (s *taskService) Update(userID, workspaceID, taskID uint, req dto.UpdateTaskRequest) (*model.Task, error) {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return nil, err
	}
	wsID, err := s.taskWorkspaceID(taskID)
	if err != nil {
		return nil, err
	}
	if wsID != workspaceID {
		return nil, repository.ErrNotFound
	}
	task, err := s.tasks.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	task.Title = req.Title
	task.Description = req.Description
	task.AssigneeID = req.AssigneeID
	task.Priority = req.Priority
	task.DueDate = req.DueDate
	if err := s.tasks.Update(task); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	if err := s.replaceTags(workspaceID, task, req.TagIDs); err != nil {
		return nil, err
	}
	s.logActivity(workspaceID, task.BoardID, &task.ID, userID, constants.ActionTaskUpdated, fmt.Sprintf("更新了任务「%s」", task.Title))
	if req.AssigneeID != nil && *req.AssigneeID != userID {
		s.notify(*req.AssigneeID, userID, constants.NotifyAssigned, fmt.Sprintf("任务「%s」已分配给你", task.Title))
	}
	return s.tasks.FindDetailed(taskID)
}

func (s *taskService) Delete(userID, workspaceID, taskID uint) error {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return err
	}
	wsID, err := s.taskWorkspaceID(taskID)
	if err != nil {
		return err
	}
	if wsID != workspaceID {
		return repository.ErrNotFound
	}
	task, _ := s.tasks.FindByID(taskID)
	if err := s.tasks.Delete(taskID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if task != nil {
		s.logActivity(workspaceID, task.BoardID, &taskID, userID, constants.ActionTaskDeleted, fmt.Sprintf("删除了任务「%s」", task.Title))
	}
	return nil
}

func (s *taskService) Move(userID, workspaceID, taskID uint, req dto.MoveTaskRequest) (*model.Task, error) {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return nil, err
	}
	wsID, err := s.taskWorkspaceID(taskID)
	if err != nil {
		return nil, err
	}
	if wsID != workspaceID {
		return nil, repository.ErrNotFound
	}
	task, err := s.tasks.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	column, err := s.columns.FindByID(req.ColumnID)
	if err != nil {
		return nil, err
	}
	if column.BoardID != task.BoardID {
		return nil, repository.ErrNotFound
	}

	if len(req.TaskIDs) > 0 {
		for i, id := range req.TaskIDs {
			t, findErr := s.tasks.FindByID(id)
			if findErr != nil {
				continue
			}
			if t.BoardID != task.BoardID {
				continue
			}
			if err := s.tasks.UpdatePosition(id, req.ColumnID, i); err != nil {
				return nil, fmt.Errorf("reorder task %d: %w", id, err)
			}
		}
	} else {
		if err := s.tasks.UpdatePosition(taskID, req.ColumnID, req.Position); err != nil {
			return nil, fmt.Errorf("move task: %w", err)
		}
	}

	oldColumnName := ""
	if task.ColumnID != req.ColumnID {
		oldColumn, _ := s.columns.FindByID(task.ColumnID)
		newColumn, _ := s.columns.FindByID(req.ColumnID)
		if oldColumn != nil && newColumn != nil {
			oldColumnName = oldColumn.Name
			s.logActivity(workspaceID, task.BoardID, &task.ID, userID, constants.ActionTaskMoved, fmt.Sprintf("将任务「%s」从「%s」移动到「%s」", task.Title, oldColumn.Name, newColumn.Name))
		}
	} else if oldColumnName == "" {
		s.logActivity(workspaceID, task.BoardID, &task.ID, userID, constants.ActionTaskMoved, fmt.Sprintf("调整了任务「%s」的排序", task.Title))
	}
	return s.tasks.FindDetailed(taskID)
}

func (s *taskService) replaceTags(workspaceID uint, task *model.Task, tagIDs []uint) error {
	for _, id := range tagIDs {
		tag, err := s.tags.FindByID(id)
		if err != nil {
			return err
		}
		if tag.WorkspaceID != workspaceID {
			return repository.ErrNotFound
		}
	}
	return s.tasks.SetTagged(task, tagIDs)
}

func (s *taskService) CreateTag(userID, workspaceID uint, req dto.CreateTagRequest) (*model.Tag, error) {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return nil, err
	}
	tag := &model.Tag{WorkspaceID: workspaceID, Name: req.Name, Color: req.Color}
	if err := s.tags.Create(tag); err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}
	return tag, nil
}

func (s *taskService) ListTags(workspaceID uint) ([]model.Tag, error) {
	return s.tags.ListByWorkspace(workspaceID)
}

func (s *taskService) DeleteTag(userID, workspaceID, tagID uint) error {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return err
	}
	tag, err := s.tags.FindByID(tagID)
	if err != nil {
		return err
	}
	if tag.WorkspaceID != workspaceID {
		return repository.ErrNotFound
	}
	if err := s.tags.Delete(tagID); err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return nil
}

func (s *taskService) CreateSubtask(userID, workspaceID, taskID uint, req dto.CreateSubtaskRequest) (*model.Subtask, error) {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return nil, err
	}
	wsID, err := s.taskWorkspaceID(taskID)
	if err != nil {
		return nil, err
	}
	if wsID != workspaceID {
		return nil, repository.ErrNotFound
	}
	subtask := &model.Subtask{TaskID: taskID, Title: req.Title, Position: 0}
	if err := s.subtasks.Create(subtask); err != nil {
		return nil, fmt.Errorf("create subtask: %w", err)
	}
	return subtask, nil
}

func (s *taskService) UpdateSubtask(userID, workspaceID, subtaskID uint, req dto.UpdateSubtaskRequest) (*model.Subtask, error) {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return nil, err
	}
	subtask, err := s.subtasks.FindByID(subtaskID)
	if err != nil {
		return nil, err
	}
	wsID, err := s.taskWorkspaceID(subtask.TaskID)
	if err != nil {
		return nil, err
	}
	if wsID != workspaceID {
		return nil, repository.ErrNotFound
	}
	subtask.Title = req.Title
	subtask.Completed = req.Completed
	if err := s.subtasks.Update(subtask); err != nil {
		return nil, fmt.Errorf("update subtask: %w", err)
	}
	return subtask, nil
}

func (s *taskService) DeleteSubtask(userID, workspaceID, subtaskID uint) error {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return err
	}
	subtask, err := s.subtasks.FindByID(subtaskID)
	if err != nil {
		return err
	}
	wsID, err := s.taskWorkspaceID(subtask.TaskID)
	if err != nil {
		return err
	}
	if wsID != workspaceID {
		return repository.ErrNotFound
	}
	if err := s.subtasks.Delete(subtaskID); err != nil {
		return fmt.Errorf("delete subtask: %w", err)
	}
	return nil
}

func (s *taskService) CreateComment(userID, workspaceID, taskID uint, req dto.CreateCommentRequest) (*model.Comment, error) {
	if err := s.requireViewer(workspaceID, userID); err != nil {
		return nil, err
	}
	wsID, err := s.taskWorkspaceID(taskID)
	if err != nil {
		return nil, err
	}
	if wsID != workspaceID {
		return nil, repository.ErrNotFound
	}
	comment := &model.Comment{TaskID: taskID, UserID: userID, Content: req.Content}
	if err := s.comments.Create(comment); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	task, _ := s.tasks.FindByID(taskID)
	if task != nil && task.AssigneeID != nil && *task.AssigneeID != userID {
		s.notify(*task.AssigneeID, userID, constants.NotifyCommented, fmt.Sprintf("任务「%s」有新评论", task.Title))
	}
	s.logActivity(workspaceID, task.BoardID, &taskID, userID, constants.ActionCommentAdded, fmt.Sprintf("评论了任务「%s」", task.Title))
	return comment, nil
}

func (s *taskService) CreateAttachment(userID, workspaceID, taskID uint, attachment *model.Attachment) (*model.Attachment, error) {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return nil, err
	}
	wsID, err := s.taskWorkspaceID(taskID)
	if err != nil {
		return nil, err
	}
	if wsID != workspaceID {
		return nil, repository.ErrNotFound
	}
	attachment.TaskID = taskID
	attachment.UploadedBy = userID
	if err := s.attachments.Create(attachment); err != nil {
		return nil, fmt.Errorf("create attachment: %w", err)
	}
	return attachment, nil
}

func (s *taskService) DeleteAttachment(userID, workspaceID, attachmentID uint) error {
	if err := s.requireEditor(workspaceID, userID); err != nil {
		return err
	}
	attachment, err := s.attachments.FindByID(attachmentID)
	if err != nil {
		return err
	}
	wsID, err := s.taskWorkspaceID(attachment.TaskID)
	if err != nil {
		return err
	}
	if wsID != workspaceID {
		return repository.ErrNotFound
	}
	if err := s.attachments.Delete(attachmentID); err != nil {
		return fmt.Errorf("delete attachment: %w", err)
	}
	return nil
}

func (s *taskService) ListBoardActivities(workspaceID, boardID uint, limit int) ([]model.ActivityLog, error) {
	if err := s.ensureBoard(workspaceID, boardID); err != nil {
		return nil, err
	}
	return s.activities.ListByBoard(boardID, limit)
}

func (s *taskService) ListWorkspaceActivities(workspaceID uint, limit int) ([]model.ActivityLog, error) {
	return s.activities.ListByWorkspace(workspaceID, limit)
}

func (s *taskService) Stats(workspaceID uint) (map[string]any, error) {
	columnCounts, err := s.tasks.CountByColumn(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("stats column count: %w", err)
	}
	assigneeCounts, err := s.tasks.CountByAssignee(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("stats assignee count: %w", err)
	}
	now := time.Now()
	from := now.AddDate(0, 0, -6)
	daily, err := s.tasks.CompletedInRange(workspaceID, from, now)
	if err != nil {
		return nil, fmt.Errorf("stats trend: %w", err)
	}
	overdue, err := s.tasks.ListOverdue(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("stats overdue: %w", err)
	}
	return map[string]any{
		"columns":       columnCounts,
		"assignees":     assigneeCounts,
		"weekly_trend":  daily,
		"overdue_tasks": overdue,
	}, nil
}

func (s *taskService) logActivity(workspaceID, boardID uint, taskID *uint, userID uint, action, detail string) {
	activity := &model.ActivityLog{
		WorkspaceID: workspaceID,
		BoardID:     boardID,
		TaskID:      taskID,
		UserID:      userID,
		Action:      action,
		Detail:      detail,
	}
	if err := s.activities.Create(activity); err != nil {
		s.logger.Error("log activity", "error", err)
	}
}

func (s *taskService) notify(userID, actorID uint, typ, content string) {
	if userID == 0 {
		return
	}
	notification := &model.Notification{
		UserID:  userID,
		ActorID: &actorID,
		Type:    typ,
		Content: content,
	}
	if err := s.notifications.Create(notification); err != nil {
		s.logger.Error("create notification", "error", err)
	}
}
