package shanglian

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/srun-soft/pay/configs"
	"github.com/srun-soft/pay/internal/app/shanglian/config"
	"github.com/srun-soft/pay/internal/app/shanglian/models"
	"github.com/srun-soft/pay/internal/database"
	"github.com/srun-soft/pay/tools/crypt"
	"net/http"
	"time"
)

func Check(c *gin.Context) {
	sd, _ := c.Get("strDes")
	res := sd.(models.StrDes)
	bp, _ := c.Get("baseParam")
	baseParam := bp.(models.Request)

	db, err := database.GetDB()
	if err != nil {
		configs.Log.Error("MySQL 连接失败", err)
		Error("MySQL 连接失败", c)
		return
	}
	var user models.User
	err = db.Where("user_name = ?", res.DeviceId).First(&user).Error
	if err != nil {
		configs.Log.Error("用户不存在", err)
		Error("用户不存在", c)
		return
	}
	// 判断状态
	if user.UserAvailable != 0 {
		configs.Log.Error("用户账号状态非正常")
		Error("用户账号状态非正常", c)
		return
	}
	// 查询产品
	cmd := database.Rdb16382.HGet(database.Rdb16382.Context(), fmt.Sprintf("hash:users:products:%d:%d", user.UserId, config.C.ProductId), "products_id")
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			configs.Log.Error("未查询到用户绑定产品")
			Error("未查询到用户绑定产品", c)
			return
		} else {
			configs.Log.Error("查询用户产品发生错误", cmd.Err())
			Error("查询用户产品发生错误", c)
		}
		return
	}

	// 组装加密数据
	baseRes := &models.BaseRes{
		ResultCode: "000000",
		ResultMsg:  "充值编码校验通过",
		TradeTime:  time.Now().Format("20060102150405"),
		DeviceId:   user.UserName,
		Detail:     "充值编码校验结果",
		DeviceInfo: models.DeviceInfo{
			NameTitle: "用户名",
			NameValue: user.UserName,
			UnitTitle: "真实姓名",
			UnitValue: user.UserRealName,
			Infos: []models.Info{
				{
					Key:      "balance",
					KeyName:  "余额",
					KeyValue: user.Balance,
				},
			},
		},
		BillFlag: 0,
	}
	jsonData, _ := json.Marshal(baseRes)
	configs.Log.Debug("响应数据", baseRes)
	strDes, _ := crypt.Encrypt(jsonData, config.C.PrivateKey)
	// 最终响应
	result := models.Response{
		BaseParam: baseParam.BaseParam,
	}
	result.StrDes = strDes
	result.SignMsg, _ = crypt.Sign(strDes, config.C.PrivateKey, config.C.Pfx)
	//
	c.JSON(http.StatusOK, result)
	return
}
