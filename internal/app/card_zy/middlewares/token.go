package middlewares

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"srun4-pay/init/common"
	"srun4-pay/init/init_card_zy"
	"srun4-pay/internal/database"
	"srun4-pay/tools"
)

type AccessToken struct {
	Token  string `json:"access_token"`
	Expire int    `json:"expires_in"`
}

// Token 获取token
func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken()
		common.Log.WithField("token", token.Token).Info()
		if token.Token == "" {
			// 获取 Token 失败
			_ = c.Error(errors.New("获取Token失败"))
			c.Next()
			return
		}
		c.Set("access_token", token)
		c.Next()
	}
}

// 获取Token
func getToken() (token AccessToken) {
	token = AccessToken{}
	var (
		res []byte
		err error
	)
	cache := database.Rdb16384
	if cache == nil {
		common.Log.WithField("Redis[16384] connect failed", "连接失败").Warn()
		return
	}

	token.Token, err = cache.Get(cache.Context(), init_card_zy.CacheAccessToken).Result()
	if err != nil && len(token.Token) > 0 {
		return
	}

	res, err = tools.Get(init_card_zy.Zy.Zhengyuan.ApiUrl)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &token)
	if err != nil {
		common.Log.WithField("获取 Token,Json解析失败", err).Error()
		return
	}
	return token
}
