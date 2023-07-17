package srun4_api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"srun4-pay/init/common"
	"srun4-pay/internal/database"
	"srun4-pay/tools"
	"strings"
)

const (
	Srun4AccessToken       = "srun4:pay:access_token"
	Srun4ApiGetAccessToken = "/api/v2/auth/get-access-token"
)

type Auth struct {
	AppId     string `gorm:"appId"`
	AppSecret string `gorm:"appSecret"`
}

func GetAccessToken() (token string) {
	var err error
	cache := database.Rdb16384
	if cache == nil {
		common.Log.WithField("Redis[16384] connect failed", "连接失败").Warn()
		return
	}
	token, err = cache.Get(cache.Context(), Srun4AccessToken).Result()
	if err != nil && len(token) > 0 {
		return
	}

	res, err = tools.Get(common.Conf.SrunConfig.)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &token)
	if err != nil {
		common.Log.WithField("获取北向接口 Token,Json解析失败", err).Error()
		return
	}
	return token
}

func getParams() {
	db := database.DB
	if db == nil {
		common.Log.WithField("MySQL connect failed", "连接失败").Warn()
		return
	}
	var auth Auth
	db.Raw("SELECT appSecret, appId FROM authorization WHERE appId = ?", "srunsoft").Scan(&auth)
	v := url.Values{}
	v.Set("appId", auth.AppId)
	v.Set("appSecret", auth.AppSecret)

	resp, err := tools.Post(fmt.Sprintf("%s%s", common.Conf.InterfaceIp, Srun4ApiGetAccessToken), v)
	if err != nil {
		common.Log.WithField("获取北向接口授权失败", err).Error()
		return
	}
	jsonData := strings.NewReader(string(resp))
	decoder := json.NewDecoder(jsonData)
	t := make(map[string]interface{})
	err = decoder.Decode(&t)

}
