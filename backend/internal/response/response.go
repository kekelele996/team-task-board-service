package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body is the unified API response envelope.
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// OK writes a successful response.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "ok", Data: data})
}

// Created writes a 201 response.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Body{Code: 0, Message: "ok", Data: data})
}

// Error writes an error response and aborts the request.
func Error(c *gin.Context, httpStatus, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Body{Code: code, Message: message, Data: nil})
}
