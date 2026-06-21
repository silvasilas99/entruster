package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"An error occurred"`
}

func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
	})
}