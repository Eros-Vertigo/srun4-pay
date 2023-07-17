package handlers

import (
	"github.com/gin-gonic/gin"
	"srun4-pay/init/common"
)

func Test(c *gin.Context) {
	c.Set("msg", "成功")
	c.Set("data", "asdasd")
	c.Next()
	common.Log.WithField("test", "handlers").Info()
}
