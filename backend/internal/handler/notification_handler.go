package handler

import (
	"strconv"

	"github.com/gbkanban/gbkanban/internal/middleware"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
)

// NotificationHandler exposes user notification endpoints.
type NotificationHandler struct {
	notifications service.NotificationService
}

// NewNotificationHandler constructs a NotificationHandler.
func NewNotificationHandler(notifications service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifications: notifications}
}

// List lists notifications for the current user.
func (h *NotificationHandler) List(c *gin.Context) {
	limit := 50
	if raw := c.Query("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}
	items, err := h.notifications.List(middleware.CurrentUserID(c), limit)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, items)
}

// UnreadCount returns the unread notification count.
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	count, err := h.notifications.UnreadCount(middleware.CurrentUserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"count": count})
}

// MarkRead marks one notification read.
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, ok := parseUintParam(c, "notification_id")
	if !ok {
		return
	}
	if err := h.notifications.MarkRead(middleware.CurrentUserID(c), id); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"read": true})
}

// MarkAllRead marks all notifications read.
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	if err := h.notifications.MarkAllRead(middleware.CurrentUserID(c)); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"read": true})
}
