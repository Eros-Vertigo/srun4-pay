package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"srun4-pay/configs"
	"time"
)

var (
	DB *gorm.DB
)

// GetDB 获取 MySQL 连接
func GetDB() (*gorm.DB, error) {
	var err error
	if DB != nil {
		return DB, nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		configs.Conf.Username,
		configs.Conf.Password,
		configs.Conf.Hostname,
		configs.Conf.Port,
		configs.Conf.Dbname)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			configs.Log,
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // 慢查询阈值
				LogLevel:                  logger.Error,           // 日志级别
				IgnoreRecordNotFoundError: true,                   // 忽略记录不存在的错误
				Colorful:                  false,                  // 禁用日志颜色
			},
		),
	})
	if err != nil {
		configs.Log.WithField("MySQL connect failed", err).Error()
		return nil, err
	}
	configs.Log.WithField("MySQL connect", "Successful")
	return DB, nil
}
