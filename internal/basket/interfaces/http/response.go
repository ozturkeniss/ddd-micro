package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendError sends an error response
func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	})
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Message: message,
	}
	
	if data != nil {
		response.Data = data
	}
	
	c.JSON(statusCode, response)
}
