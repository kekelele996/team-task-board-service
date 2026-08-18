package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gbkanban/gbkanban/internal/config"
	"github.com/gbkanban/gbkanban/internal/handler"
	"github.com/gbkanban/gbkanban/internal/repository"
	"github.com/gbkanban/gbkanban/internal/router"
	"github.com/gbkanban/gbkanban/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load configuration", "error", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		logger.Error("create upload directory", "path", cfg.UploadDir, "error", err)
		os.Exit(1)
	}

	db, err := repository.NewDB(cfg.DSN())
	if err != nil {
		logger.Error("initialize database", "error", err)
		os.Exit(1)
	}

	users := repository.NewUserRepo(db)
	workspaces := repository.NewWorkspaceRepo(db)
	boards := repository.NewBoardRepo(db)
	columns := repository.NewColumnRepo(db)
	tasks := repository.NewTaskRepo(db)
	tags := repository.NewTagRepo(db)
	subtasks := repository.NewSubtaskRepo(db)
	comments := repository.NewCommentRepo(db)
	attachments := repository.NewAttachmentRepo(db)
	activities := repository.NewActivityRepo(db)
	notifications := repository.NewNotificationRepo(db)

	workspaceService := service.NewWorkspaceService(workspaces, users, logger)
	boardService := service.NewBoardService(boards, columns, workspaces, workspaceService, logger)
	taskService := service.NewTaskService(tasks, boards, columns, tags, subtasks, comments, attachments, activities, notifications, workspaceService, logger)
	notificationService := service.NewNotificationService(notifications, logger)
	authService := service.NewAuthService(users, cfg.JWTSecret, cfg.TokenTTLHours, logger)

	app := router.NewRouter(
		authService,
		workspaceService,
		handler.NewAuthHandler(authService, users),
		handler.NewWorkspaceHandler(workspaceService),
		handler.NewBoardHandler(boardService),
		handler.NewTaskHandler(taskService, cfg.UploadDir),
		handler.NewNotificationHandler(notificationService),
		cfg.UploadDir,
		cfg.AllowedOrigin,
		logger,
	).Engine()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           app,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown server", "error", err)
		}
	}()

	logger.Info("server started", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("serve HTTP", "error", err)
		os.Exit(1)
	}
}
