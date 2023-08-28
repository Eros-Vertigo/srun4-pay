package models

import (
	"github.com/shopspring/decimal"
)

type Alipay struct {
	ID         int             `gorm:"column:id;primaryKey;autoIncrement"`
	UserName   string          `gorm:"column:user_name;size:32"`
	OutTradeNo string          `gorm:"column:out_trade_no;size:32;unique"`
	Money      decimal.Decimal `gorm:"column:money"`
	Type       int             `gorm:"column:type"`
	BuyTime    int             `gorm:"column:buy_time"`
	Status     string          `gorm:"column:status;size:1;default:'0'"`
	Payment    string          `gorm:"column:payment;size:32"`
	TradeNo    string          `gorm:"column:trade_no;size:32"`
	PayType    string          `gorm:"column:pay_type;size:5;default:'0'"`
	Remark     string          `gorm:"column:remark;size:256"`
}

// TableName 表名
func (Alipay) TableName() string {
	return "alipay"
}
