package router

import (
	"log/slog"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/handler"
	"github.com/gbkanban/gbkanban/internal/middleware"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
)

// Router wires every HTTP endpoint and middleware.
type Router struct {
	engine              *gin.Engine
	auth                service.AuthService
	workspaceService    service.WorkspaceService
	authHandler         *handler.AuthHandler
	workspaceHandler    *handler.WorkspaceHandler
	boardHandler        *handler.BoardHandler
	taskHandler         *handler.TaskHandler
	notificationHandler *handler.NotificationHandler
}

// NewRouter constructs the application router.
func NewRouter(
	authService service.AuthService,
	workspaceService service.WorkspaceService,
	authHandler *handler.AuthHandler,
	workspaceHandler *handler.WorkspaceHandler,
	boardHandler *handler.BoardHandler,
	taskHandler *handler.TaskHandler,
	notificationHandler *handler.NotificationHandler,
	uploadDir string,
	allowedOrigin string,
	logger *slog.Logger,
) *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(middleware.Recovery(logger), middleware.RequestLog(logger), middleware.CORS(allowedOrigin))
	engine.Static("/uploads", uploadDir)
	engine.GET("/healthz", handler.Healthz)
	engine.GET("/health", handler.Health)

	api := engine.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		authorized := api.Group("")
		authorized.Use(middleware.Auth(authService))
		{
			authorized.GET("/me", authHandler.Me)

			authorized.GET("/workspaces", workspaceHandler.List)
			authorized.POST("/workspaces", workspaceHandler.Create)

			workspace := authorized.Group("/workspaces/:workspace_id")
			{
				workspace.GET("", workspaceHandler.Get)
				workspace.PUT("", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin), workspaceHandler.Update)
				workspace.DELETE("", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin), workspaceHandler.Delete)

				workspace.GET("/members", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), workspaceHandler.ListMembers)
				workspace.POST("/members", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin), workspaceHandler.AddMember)
				workspace.PUT("/members/:user_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin), workspaceHandler.UpdateMember)
				workspace.DELETE("/members/:user_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin), workspaceHandler.RemoveMember)

				workspace.GET("/boards", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), boardHandler.List)
				workspace.POST("/boards", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), boardHandler.Create)
				workspace.GET("/boards/:board_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), boardHandler.Get)
				workspace.PUT("/boards/:board_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), boardHandler.Update)
				workspace.DELETE("/boards/:board_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin), boardHandler.Delete)

				workspace.GET("/boards/:board_id/columns", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), boardHandler.ListColumns)
				workspace.POST("/boards/:board_id/columns", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), boardHandler.CreateColumn)
				workspace.PUT("/columns/:column_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), boardHandler.UpdateColumn)
				workspace.DELETE("/columns/:column_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), boardHandler.DeleteColumn)

				workspace.GET("/boards/:board_id/tasks", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), taskHandler.ListTasks)
				workspace.POST("/boards/:board_id/tasks", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.CreateTask)
				workspace.GET("/tasks/:task_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), taskHandler.GetTask)
				workspace.PUT("/tasks/:task_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.UpdateTask)
				workspace.DELETE("/tasks/:task_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.DeleteTask)
				workspace.PATCH("/tasks/:task_id/move", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.MoveTask)

				workspace.GET("/tags", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), taskHandler.ListTags)
				workspace.POST("/tags", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.CreateTag)
				workspace.DELETE("/tags/:tag_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.DeleteTag)

				workspace.POST("/tasks/:task_id/subtasks", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.CreateSubtask)
				workspace.PUT("/subtasks/:subtask_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.UpdateSubtask)
				workspace.DELETE("/subtasks/:subtask_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.DeleteSubtask)

				workspace.POST("/tasks/:task_id/comments", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), taskHandler.CreateComment)
				workspace.POST("/tasks/:task_id/attachments", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.CreateAttachment)
				workspace.DELETE("/attachments/:attachment_id", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor), taskHandler.DeleteAttachment)

				workspace.GET("/activities", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), taskHandler.WorkspaceActivities)
				workspace.GET("/boards/:board_id/activities", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), taskHandler.BoardActivities)
				workspace.GET("/stats", middleware.RequireWorkspaceRole(workspaceService, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer), taskHandler.Stats)
			}

			authorized.GET("/notifications", notificationHandler.List)
			authorized.GET("/notifications/unread-count", notificationHandler.UnreadCount)
			authorized.PUT("/notifications/:notification_id/read", notificationHandler.MarkRead)
			authorized.PUT("/notifications/read-all", notificationHandler.MarkAllRead)
		}
	}

	return &Router{
		engine:              engine,
		auth:                authService,
		workspaceService:    workspaceService,
		authHandler:         authHandler,
		workspaceHandler:    workspaceHandler,
		boardHandler:        boardHandler,
		taskHandler:         taskHandler,
		notificationHandler: notificationHandler,
	}
}

// Engine returns the underlying gin engine.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
