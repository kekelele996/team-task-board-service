package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/repository"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Error(c, http.StatusBadRequest, constants.CodeBadRequest, "invalid json body: "+err.Error())
		return false
	}
	if err := validate.Struct(obj); err != nil {
		response.Error(c, http.StatusBadRequest, constants.CodeBadRequest, formatValidationError(err))
		return false
	}
	return true
}

func bindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		response.Error(c, http.StatusBadRequest, constants.CodeBadRequest, "invalid query params: "+err.Error())
		return false
	}
	return true
}

func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		return validationErrors[0].Error()
	}
	return err.Error()
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, constants.CodeUnauthorized, "invalid email or password")
	case errors.Is(err, service.ErrEmailTaken):
		response.Error(c, http.StatusConflict, constants.CodeConflict, "email already registered")
	case errors.Is(err, service.ErrUsernameTaken):
		response.Error(c, http.StatusConflict, constants.CodeConflict, "username already taken")
	case errors.Is(err, service.ErrForbidden):
		response.Error(c, http.StatusForbidden, constants.CodeForbidden, "forbidden")
	case errors.Is(err, service.ErrInvalidMember):
		response.Error(c, http.StatusNotFound, constants.CodeNotFound, "member not found")
	case errors.Is(err, repository.ErrNotFound):
		response.Error(c, http.StatusNotFound, constants.CodeNotFound, "resource not found")
	case errors.Is(err, repository.ErrConflict):
		response.Error(c, http.StatusConflict, constants.CodeConflict, "resource already exists")
	default:
		response.Error(c, http.StatusInternalServerError, constants.CodeInternal, "internal server error")
	}
}

func parseUintParam(c *gin.Context, key string) (uint, bool) {
	raw := c.Param(key)
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		response.Error(c, http.StatusBadRequest, constants.CodeBadRequest, "invalid "+key)
		return 0, false
	}
	return uint(value), true
}
