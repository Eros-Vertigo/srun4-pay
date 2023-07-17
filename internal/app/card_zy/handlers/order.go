package handlers

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"srun4-pay/init/init_card_zy"
	"srun4-pay/tools"
	"strings"
	"time"
)

type PreOrderParams struct {
	UniqueId   string `form:"uniqueId"`
	QueryType  string `form:"queryType"`
	EWalletNum string `form:"eWalletNum"`
	MonTrans   string `form:"monTrans"`
	DealerNum  string `form:"dealerNum"`
	DealTime   string `form:"dealTime"`
	Sign       string `form:"sign"`
}

// Preorder 预下单
func Preorder(c *gin.Context) {
	accessToken, _ := c.Get("access_token")
	// 组装请求参数
	p := &PreOrderParams{
		UniqueId:   c.Query("uniqueId"),
		QueryType:  "2",
		EWalletNum: "1",
		MonTrans:   c.Query("monTrans"),
		DealerNum:  init_card_zy.Zy.Zhengyuan.MerchantNo,
		DealTime:   time.Now().Format("2006-01-02 15:04:05"),
	}

	temp := tools.GenerateSignature(p)
	p.Sign = fmt.Sprintf("%X", md5.Sum([]byte(fmt.Sprintf("%s&key=%s", temp, init_card_zy.Zy.Zhengyuan.AppId))))
	v := url.Values{}
	v.Set("optType", "5")
	v.Set("uniqueId", p.UniqueId)
	v.Set("queryType", p.QueryType)
	v.Set("eWalletNum", p.EWalletNum)
	v.Set("monTrans", p.MonTrans)
	v.Set("dealerNum", p.DealerNum)
	v.Set("dealTime", p.DealTime)
	v.Set("sign", p.Sign)

	resp, err := tools.Post(fmt.Sprintf("%s%s?access_token=%s", init_card_zy.Zy.Zhengyuan.ApiUrl, init_card_zy.MethodPreOrder, accessToken), v)
	if err != nil {
		_ = c.Error(errors.New(err.Error()))
		c.Next()
		return
	}

	jsonData := strings.NewReader(string(resp))
	decoder := json.NewDecoder(jsonData)
	t := make(map[string]interface{})
	err = decoder.Decode(&t)
	if err != nil {
		_ = c.Error(errors.New(err.Error()))
		c.Next()
		return
	}
	// 生成订单号
	c.Set("trade_no", tools.GenerateOrderNumber())
	c.Set("preorder", p)
	c.Set("recId", t["recId"])
	c.Next()
	return
}
