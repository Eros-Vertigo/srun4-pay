package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"srun4-pay/configs"
	"srun4-pay/internal/app/shanglian/config"
	"srun4-pay/internal/app/shanglian/models"
	crypt2 "srun4-pay/tools/crypt"
)

// 验证

func VerifySignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		configs.Log.Debug("verify signature")
		var req models.Request
		if err := c.BindJSON(&req); err != nil {
			configs.Log.Error("请求参数解析错误")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			c.Done()
			return
		}

		// 在这里进行验签和解密的操作
		// 根据需要对解密后的数据进行处理
		// 使用crypt包中的签名验证方法进行验签
		valid, _ := crypt2.Verify(req.StrDes, req.SignMsg, config.C.Cert)
		if !valid {
			configs.Log.Error("签名验证失败")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Signature verification failed"})
			return
		}

		// 使用crypt包中的解密方法进行解密
		decrypted, err := crypt2.Decrypt(req.StrDes, config.C.PrivateKey)
		if err != nil {
			configs.Log.Error("解密失败:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
			return
		}

		var data models.StrDes
		err = json.Unmarshal(decrypted, &data)
		if err != nil {
			configs.Log.Error("json 解析失败:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Json parse error"})
			return
		}
		c.Set("strDes", data)
		c.Set("baseParam", req)
		// 将解密后的参数重新绑定到上下文中，以便后续处理函数使用
		c.Next()
	}
}
