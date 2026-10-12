package handlers

import "github.com/gin-gonic/gin"

func (h BaseHandler) GetQueryParam(c *gin.Context, param, defaultValue string) string {
	return c.DefaultQuery(param, defaultValue)
}
