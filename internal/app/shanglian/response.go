package shanglian

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/srun-soft/pay/configs"
	"github.com/srun-soft/pay/internal/app/shanglian/config"
	"github.com/srun-soft/pay/internal/app/shanglian/models"
	crypt2 "github.com/srun-soft/pay/tools/crypt"
	"net/http"
	"time"
)

func Error(msg string, c *gin.Context) {
	bp, _ := c.Get("baseParam")
	baseParam := bp.(models.Request)
	errRes := &models.BaseRes{
		ResultCode: "E99999",
		ResultMsg:  msg,
		TradeTime:  time.Now().Format("20060102150405"),
		DeviceId:   "",
		Detail:     "错误响应",
		DeviceInfo: models.DeviceInfo{
			NameTitle: "用户名",
			NameValue: "",
			UnitTitle: "真实姓名",
			UnitValue: "",
			Infos: []models.Info{
				{
					Key:      "balance",
					KeyName:  "余额",
					KeyValue: "",
				},
			},
		},
		BillFlag: 0,
	}
	jsonData, _ := json.Marshal(errRes)
	configs.Log.WithField("返回错误响应", errRes.ResultMsg).Info()
	encrypted, _ := crypt2.Encrypt(jsonData, config.C.PrivateKey)
	signMsg, _ := crypt2.Sign(encrypted, config.C.PrivateKey, config.C.Pfx)

	c.JSON(http.StatusOK, models.Response{
		BaseParam: baseParam.BaseParam,
		StrDes:    encrypted,
		SignMsg:   signMsg,
	})
}
