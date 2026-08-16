package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/repository"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
)

// RequireWorkspaceRole checks that the authenticated user has one of the
// allowed roles in the workspace identified by the `workspace_id` path param.
func RequireWorkspaceRole(workspaceService service.WorkspaceService, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.Param("workspace_id")
		if raw == "" {
			raw = c.Query("workspace_id")
		}
		workspaceID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, constants.CodeBadRequest, "invalid workspace_id")
			return
		}
		if err := workspaceService.RequireRole(uint(workspaceID), CurrentUserID(c), roles...); err != nil {
			if errors.Is(err, service.ErrForbidden) {
				response.Error(c, http.StatusForbidden, constants.CodeForbidden, "you do not have permission for this workspace")
				return
			}
			if errors.Is(err, repository.ErrNotFound) {
				response.Error(c, http.StatusNotFound, constants.CodeNotFound, "workspace membership not found")
				return
			}
			response.Error(c, http.StatusInternalServerError, constants.CodeInternal, "internal server error")
			return
		}
		c.Next()
	}
}
