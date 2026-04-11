package utils

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func Error(c *gin.Context, message string, code int, err ...error) {
	var errMsg string

	if len(err) > 0 && err[0] != nil {
		errMsg = err[0].Error()
	}

	c.JSON(code, ErrorResponse{
		Success: false,
		Message: message,
		Error:   errMsg,
	})
}
