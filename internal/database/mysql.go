package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"srun4-pay/init/common"
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
		common.Conf.Username,
		common.Conf.Password,
		common.Conf.Hostname,
		common.Conf.Port,
		common.Conf.Dbname)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			common.Log,
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // 慢查询阈值
				LogLevel:                  logger.Error,           // 日志级别
				IgnoreRecordNotFoundError: true,                   // 忽略记录不存在的错误
				Colorful:                  false,                  // 禁用日志颜色
			},
		),
	})
	if err != nil {
		common.Log.WithField("MySQL connect failed", err).Error()
		return nil, err
	}
	common.Log.WithField("MySQL connect", "Successful")
	return DB, nil
}
