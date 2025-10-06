package handlers

import (
	"github.com/gin-gonic/gin"
)

func sendErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
