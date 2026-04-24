package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
		Data:    data,
	})
}
