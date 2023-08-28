package shanglian

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"srun4-pay/configs"
	"srun4-pay/internal/app/shanglian/config"
	"srun4-pay/internal/app/shanglian/models"
	"srun4-pay/internal/database"
	"srun4-pay/tools/crypt"
	"strconv"
	"time"
)

func Notify(c *gin.Context) {
	bp, _ := c.Get("baseParam")
	baseParam := bp.(models.Request)
	sd, _ := c.Get("strDes")
	strDes := sd.(models.StrDes)
	// 创建alipay订单
	alipay := models.Alipay{
		UserName:   strDes.DeviceId,
		OutTradeNo: baseParam.TradeNo,
		Money:      strDes.Amount,
		Type:       config.C.ProductId,
		BuyTime:    int(time.Now().Unix()),
		Status:     strconv.Itoa(1),
		Payment:    strDes.DeviceId,
		TradeNo:    baseParam.TradeNo,
		PayType:    strconv.Itoa(config.C.PayType),
		Remark:     "商联缴费",
	}
	db, err := database.GetDB()
	if err != nil {
		configs.Log.Error("MySQL 连接失败", err)
		Error("MySQL 连接失败", c)
		return
	}
	save := db.Create(&alipay)
	if save.Error != nil {
		configs.Log.Error("订单创建失败", save.Error.Error())
		Error("订单创建失败", c)
		return
	}

	code := recharge(&alipay)
	if !code {
		Error("订单处理失败", c)
		return
	}
	res := models.NotifyRes{
		BaseRes: models.BaseRes{
			ResultCode: "000000",
			ResultMsg:  "充值响应处理完成",
			TradeTime:  time.Now().Format("20060102150405"),
			DeviceId:   strDes.DeviceId,
			Detail:     "充值结果通知回执",
		},
		ThirdResult: "S",
	}
	jsonData, _ := json.Marshal(res)
	encryted, _ := crypt.Encrypt(jsonData, config.C.PrivateKey)
	// 最终响应
	result := models.Response{
		BaseParam: baseParam.BaseParam,
	}
	result.StrDes = encryted
	result.SignMsg, _ = crypt.Sign(encryted, config.C.PrivateKey, config.C.Pfx)
	//
	c.JSON(http.StatusOK, result)
	return
}

// 调用v2接口 进行产品充值订单
func recharge(a *models.Alipay) bool {
	data := make(map[string]interface{})
	data["user_name"] = a.UserName
	data["pay_type_id"] = config.C.PayTypeId
	data["order_no"] = a.TradeNo
	data["pay_num"] = fmt.Sprintf("%v", a.Money)
	data["amount"] = fmt.Sprintf("%v", a.Money)
	data["product"] = fmt.Sprintf("%v", a.Type)

	configs.Log.WithField("组装参数请求北向接口", data).Debug()

	res, err := config.API.ProductRecharge(data)
	if err != nil {
		configs.Log.Error("充值失败", err)
		return false
	}
	configs.Log.Info("订单处理回执", res)
	if res.Code == 0 {
		return true
	}
	return false
}
