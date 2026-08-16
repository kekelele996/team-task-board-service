package middleware

import (
	"strings"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
)

const userIDKey = "user_id"

// Auth validates the Bearer JWT and stores the user id in the context.
func Auth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, 401, constants.CodeUnauthorized, "missing authorization header")
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, 401, constants.CodeUnauthorized, "invalid authorization header")
			return
		}
		userID, err := authService.ParseToken(parts[1])
		if err != nil {
			response.Error(c, 401, constants.CodeUnauthorized, "invalid or expired token")
			return
		}
		c.Set(userIDKey, userID)
		c.Next()
	}
}

// CurrentUserID returns the authenticated user id stored by Auth.
func CurrentUserID(c *gin.Context) uint {
	value, ok := c.Get(userIDKey)
	if !ok {
		return 0
	}
	id, ok := value.(uint)
	if !ok {
		return 0
	}
	return id
}
