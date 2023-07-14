package middlewares

import (
	"github.com/gin-gonic/gin"
	"srun4-pay/init/common"
	"srun4-pay/init/init_card_zy"
	"srun4-pay/internal/database"
	"srun4-pay/internal/http"
)

// Token 获取token
func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		getToken()
		c.Set("access_token", "")
		c.Next()
	}
}

func getToken() string {
	cache := database.Rdb16384
	if cache == nil {
		common.Log.WithField("Redis[16384] connect failed", "连接失败").Warn()
	}
	t, err := cache.Get(cache.Context(), init_card_zy.CacheAccessToken).Result()
	if err != nil && len(t) > 0 {
		return t
	}

	b, err := http.Get(init_card_zy.Zy.Zhengyuan.ApiUrl)
	if err != nil {
		return ""
	}
	common.Log.WithField("获取access_token", b)
	return string(b)
}
