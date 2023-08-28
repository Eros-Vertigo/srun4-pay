package models

import (
	"github.com/shopspring/decimal"
)

// 结构体

type BaseParam struct {
	Version string      `json:"version" default:"V1.0"`
	Charset string      `json:"charset"`
	TradeNo string      `json:"tradeNo"`
	RechId  interface{} `json:"rechId"`
	IntId   string      `json:"intId"`
	MerId   int         `json:"merId"`
}

type Request struct {
	BaseParam
	StrDes  string `json:"strDes"`
	SignMsg string `json:"signMsg"`
}

type StrDes struct {
	ThirdMerId string          `json:"thirdMerId"`
	DeviceId   string          `json:"deviceId"`
	TradeTime  any             `json:"tradeTime"`
	Detail     string          `json:"detail"`
	Amount     decimal.Decimal `json:"amount"`
	Bills
}

type Bills struct {
}

type Response struct {
	BaseParam
	StrDes  string `json:"strDes"`
	SignMsg string `json:"signMsg"`
}

type BaseRes struct {
	ResultCode string     `json:"resultCode"`
	ResultMsg  string     `json:"resultMsg"`
	TradeTime  string     `json:"tradeTime"`
	DeviceId   string     `json:"deviceId"`
	Detail     string     `json:"detail"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

type NotifyRes struct {
	BaseRes
	ThirdSerioNo string `json:"thirdSerioNo"`
	ThirdResult  string `json:"thirdResult"`
}

type DeviceInfo struct {
	NameTitle string `json:"nameTitle"`
	NameValue string `json:"nameValue"`
	UnitTitle string `json:"unitTitle"`
	UnitValue string `json:"unitValue"`
	Infos     []Info `json:"infos"`
}

type Info struct {
	Key      string `json:"key"`
	KeyName  string `json:"keyName"`
	KeyValue string `json:"keyValue"`
}

type User struct {
	UserId        uint   `gorm:"user_id"`
	UserName      string `gorm:"user_name"`
	UserRealName  string `gorm:"user_real_name"`
	Balance       string `gorm:"type:decimal(10,2)"`
	UserAvailable uint   `gorm:"user_available"`
}

type IPBindingToken struct {
	Token string `gorm:"column:token"`
}
