package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// A standard response struct for API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Function to send a success response
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Function to send an error response
func ErrorResponse(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, APIResponse{
		Success: false,
		Message: message,
	})
}
