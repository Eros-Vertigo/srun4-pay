package middlewares

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"srun4-pay/init/common"
	"srun4-pay/init/init_card_zy"
	"srun4-pay/internal/app/card_zy/handlers"
	"srun4-pay/tools"
	"strings"
	"time"
)

type OrderHandleParams struct {
	TransRecId string `json:"transRecId"`
	DealTime   string `json:"dealTime"`
	ProofNum   string `json:"proofNum"`
	RecDate    string `json:"recDate"`
	PayType    string `json:"payType"`
	Sign       string `json:"sign"`
}

func OrderHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Errors.Last() != nil {
			c.Next()
			return
		}

		if recId, exist := c.Get("recId"); exist {
			tradeNo, _ := c.Get("trade_no")
			accessToken, _ := c.Get("access_token")
			p := OrderHandleParams{
				TransRecId: recId.(string),
				DealTime:   time.Now().Format("2006-01-02 15:04:05"),
				ProofNum:   tradeNo.(string),
				RecDate:    time.Now().Format("2006-01-02"),
				PayType:    "2",
			}
			temp := tools.GenerateSignature(p)
			p.Sign = fmt.Sprintf("%X", md5.Sum([]byte(fmt.Sprintf("%s&key=%s", temp, init_card_zy.Zy.Zhengyuan.AppId))))

			v := url.Values{}
			v.Set("transRecId", p.TransRecId)
			v.Set("dealTime", p.DealTime)
			v.Set("proofNum", p.ProofNum)
			v.Set("recDate", p.RecDate)
			v.Set("payType", p.PayType)
			v.Set("sign", p.Sign)
			resp, err := tools.Post(fmt.Sprintf("%s%s?access_token=%s", init_card_zy.Zy.Zhengyuan.ApiUrl, init_card_zy.MethodOrderHandle, accessToken), v)
			if err != nil {
				_ = c.Error(errors.New(err.Error()))
				c.Next()
				return
			}

			jsonData := strings.NewReader(string(resp))
			decoder := json.NewDecoder(jsonData)
			var t map[string]interface{}
			err = decoder.Decode(&t)

			if _, ok := t["code"]; !ok {
				_ = c.Error(errors.New("响应参数未获取到code码"))
				c.Next()
				return
			}
			preorder, _ := c.Get("preorder")

			// 调用北向接口进行充值
			err = Recharge(preorder.(handlers.PreOrderParams))
			if err != nil {
				_ = c.Error(errors.New(err.Error()))
				c.Next()
				return
			}
			c.JSON(http.StatusOK, t)
			return
		}
	}
}

func Recharge(pre handlers.PreOrderParams) error {
	v := url.Values{}
	v.Set("user_name", "")
	v.Set("pay_type_id", "2")
	v.Set("order_no", pre.MonTrans)
	v.Set("pay_num", pre.MonTrans)

	common.Log.WithField("一卡通订单处理成功,调用北向接口进行充值", v).Info()
	// 调用北向接口进行充值
	resp, err := tools.Post("", v)
	if err != nil {
		return err
	}
	common.Log.WithField("北向接口充值响应", resp).Info()

	jsonData := strings.NewReader(string(resp))
	decoder := json.NewDecoder(jsonData)
	var t map[string]interface{}
	err = decoder.Decode(&t)
	if code, _ := t["code"]; code != 0 {
		return errors.New(fmt.Sprintf("充值失败:%s", t["message"]))
	}
	return nil
}
