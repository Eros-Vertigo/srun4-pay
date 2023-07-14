package middlewares

import "github.com/gin-gonic/gin"

// Token 获取token
func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("access_token", "")
		c.Next()
	}
}
