package shanglian

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"srun4-pay/configs"
	"time"
)

// Load 加载Yaml配置文件
func (c *Config) Load() {
	err := yaml.Unmarshal(configs.PayConf, &c)
	if err != nil {
		configs.Log.WithField("yaml 解析失败,err:%s", err).Fatal()
	}
	configs.Log.WithField(fmt.Sprintf("支付【%s】配置", configs.PayYaml), "加载成功").Info()
	// 日志目录
	err = os.MkdirAll(c.LogPath, os.ModePerm)
	if err != nil {
		configs.Log.WithField(fmt.Sprintf("Failed to create log directory:%s", c.LogPath), err).Fatal()
	}
	// 设置日志文件名
	logFileName := filepath.Join(c.LogPath, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		configs.Log.WithField("Failed to open log file:%s", err).Fatal()
	}

	if *configs.Mode == "prod" {
		mw := io.MultiWriter(os.Stdout, logFile)
		configs.Log.SetOutput(mw)
	} else {
		configs.Log.SetOutput(os.Stdout)
	}
}