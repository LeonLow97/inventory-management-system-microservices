package handlers

import "github.com/gin-gonic/gin"

func (h BaseHandler) ResponseJSON(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, gin.H{"data": data})
}

func (h BaseHandler) ResponseNoContent(c *gin.Context, statusCode int) {
	c.Status(statusCode)
}
