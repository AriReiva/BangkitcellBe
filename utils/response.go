package utils

import "github.com/gin-gonic/gin"

func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"success": true, "data": data})
}

func RespondMessage(c *gin.Context, message string) {
	c.JSON(200, gin.H{"success": true, "message": message})
}

func RespondError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"success": false, "error": err.Error()})
}