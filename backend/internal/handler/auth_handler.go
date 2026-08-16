package handler

import (
	"net/http"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/repository"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler exposes registration and login endpoints.
type AuthHandler struct {
	auth  service.AuthService
	users repository.UserRepo
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(auth service.AuthService, users repository.UserRepo) *AuthHandler {
	return &AuthHandler{auth: auth, users: users}
}

// Register creates a new user account.
// @Summary Register
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "registration payload"
// @Success 201 {object} response.Body
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}
	token, err := h.auth.Register(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, token)
}

// Login authenticates a user.
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "login payload"
// @Success 200 {object} response.Body
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	token, err := h.auth.Login(req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, token)
}

// Me returns the authenticated user profile.
// @Summary Me
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Router /api/v1/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID := currentUserID(c)
	user, err := h.users.FindByID(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, dto.UserDTO{ID: user.ID, Username: user.Username, Email: user.Email})
}

func currentUserID(c *gin.Context) uint {
	if id, ok := c.Get("user_id"); ok {
		if uid, ok := id.(uint); ok {
			return uid
		}
	}
	response.Error(c, http.StatusUnauthorized, constants.CodeUnauthorized, "unauthorized")
	return 0
}
