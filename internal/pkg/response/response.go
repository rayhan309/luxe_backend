package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard API response envelope
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Success sends a 200 OK response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Paginated sends a paginated list response
func Paginated(c *gin.Context, message string, data interface{}, page, limit int, total int64) {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// BadRequest sends a 400 error response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error:   message,
	})
}

// Unauthorized sends a 401 response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Error:   message,
	})
}

// Forbidden sends a 403 response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Error:   message,
	})
}

// NotFound sends a 404 response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Error:   message,
	})
}

// Conflict sends a 409 response
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, APIResponse{
		Success: false,
		Error:   message,
	})
}

// InternalError sends a 500 response
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Error:   message,
	})
}

// ValidationError sends a 422 with field errors
func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Error:   errors,
	})
}
