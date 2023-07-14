package handlers

import (
	"github.com/gin-gonic/gin"
	"srun4-pay/init/common"
)

func Test(c *gin.Context) {
	common.Log.WithField("test", "1").Info()
}
