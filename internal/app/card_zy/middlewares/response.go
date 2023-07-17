package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ErrorRes struct {
	BaseRes
}

type SuccessRes struct {
	BaseRes
	Data any `json:"data"`
}

// ErrMiddleware 错误响应
func ErrMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 捕获err
		err := c.Errors.Last()
		if err != nil {
			res := &ErrorRes{
				BaseRes: BaseRes{
					Code: 0,
					Msg:  err.Err.Error(),
				},
			}
			c.JSON(http.StatusBadRequest, res)
			c.Abort()
		}
	}
}

// SuccessMiddleware 成功响应
func SuccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if msg, exists := c.Get("msg"); exists {
			data, _ := c.Get("data")
			res := &SuccessRes{
				BaseRes: BaseRes{
					Code: 200,
					Msg:  msg.(string),
				},
				Data: data,
			}
			c.JSON(http.StatusOK, res)
			c.Abort()
		}
	}
}
