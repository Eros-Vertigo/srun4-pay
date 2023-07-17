package tools

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderNumber() string {
	now := time.Now()
	sec := now.Unix()
	nano := now.UnixNano()

	// 获取微秒数
	usec := nano / 1000

	// 格式化订单号
	orderNo := fmt.Sprintf("%s%d%04d%d", now.Format("20060102150405"), sec, usec, rand.Intn(900)+100)
	return orderNo
}
