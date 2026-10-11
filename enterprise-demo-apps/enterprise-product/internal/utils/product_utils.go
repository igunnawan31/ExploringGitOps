package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, code int, errCode string, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   errCode,
		"message": message,
	})
}
