package shanglian

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/srun-soft/api-sdk/tools"
	"github.com/srun-soft/pay/configs"
	"github.com/srun-soft/pay/internal/app/shanglian/config"
	"github.com/srun-soft/pay/internal/app/shanglian/models"
	"github.com/srun-soft/pay/internal/database"
	"github.com/srun-soft/pay/tools/crypt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const AccessToken = "access_token:shanglian"

func Notify(c *gin.Context) {
	bp, _ := c.Get("baseParam")
	baseParam := bp.(models.Request)
	sd, _ := c.Get("strDes")
	strDes := sd.(models.StrDes)

	// find user
	userID := database.Rdb16382.Get(database.Rdb16382.Context(), fmt.Sprintf("key:users:user_name:%s", strDes.DeviceId)).Val()
	// 查询产品
	cmd := database.Rdb16382.HGet(database.Rdb16382.Context(), fmt.Sprintf("hash:users:%s", userID), "products_id").Val()
	ProductID, err := strconv.Atoi(cmd)
	if err != nil {
		configs.Log.Error("获取产品ID 失败", err)
		Error("获取产品ID 失败", c)
		return
	}

	// 创建alipay订单
	alipay := models.Alipay{
		UserName:   strDes.DeviceId,
		OutTradeNo: baseParam.TradeNo,
		Money:      strDes.Amount,
		Type:       ProductID,
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

	token, err := getToken(db)
	if err != nil {
		Error("获取token失败", c)
		return
	}
	code := rechargeV1(&alipay, token)
	//code := recharge(&alipay)
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
			BillFlag:   0,
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

func getToken(db *gorm.DB) (string, error) {
	token, err := database.Rdb16384.Get(database.Rdb16384.Context(), AccessToken).Result()
	if err == nil && len(token) > 0 {
		return token, nil
	}
	var ipBindingToken models.IPBindingToken
	err = db.Table("ipBindingToken").Select("token").Where("ip like?", "127.0.0.1%").First(&ipBindingToken).Error
	if err != nil {
		configs.Log.Error("查询v1token失败", err)
		return "", err
	}

	token = ipBindingToken.Token
	err = database.Rdb16384.Set(database.Rdb16384.Context(), AccessToken, token, 3600).Err()
	if err != nil {
		configs.Log.Error("缓存v1token失败", err)
		return "", err
	}
	return token, nil
}

func rechargeV1(a *models.Alipay, token string) bool {
	data := url.Values{}
	data.Set("user_name", a.UserName)
	data.Set("pay_type_id", strconv.Itoa(config.C.PayTypeId))
	data.Set("order_no", a.TradeNo)
	data.Set("pay_num", fmt.Sprintf("%v", a.Money))
	data.Set("amount", fmt.Sprintf("%v", a.Money))
	data.Set("product", fmt.Sprintf("%v", a.Type))
	data.Set("access_token", token)

	urlPath := fmt.Sprintf("%s://%s:8001/api/v1/product/recharge", config.C.Scheme, config.C.InterfaceIP)

	configs.Log.WithField("组装参数请求北向接口", data).Debug()
	res, err := PostRequestWithoutCert(urlPath, data)
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

func PostRequestWithoutCert(url string, body url.Values) (tools.SrunResponse, error) {
	var sr tools.SrunResponse
	// 创建一个忽略证书验证的 HTTP 客户端
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// 发起 POST 请求
	resp, err := client.PostForm(url, body)
	if err != nil {
		log.Println("POST 请求失败:", err)
		return sr, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			configs.Log.Error("body close failed:", err)
		}
	}(resp.Body)

	// 读取响应内容
	err = json.NewDecoder(resp.Body).Decode(&sr)
	if err != nil {
		fmt.Println("解析响应失败:", err)
		return sr, fmt.Errorf("POST 解析响应失败: %s", resp.Status)
	}

	return sr, nil
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
