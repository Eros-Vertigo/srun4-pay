package shanglian

import "github.com/gin-gonic/gin"

func Notify(c *gin.Context) {
	c.JSON(200, "notify")
	return
}
