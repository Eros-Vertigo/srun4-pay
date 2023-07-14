package init_card_zy

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"srun4-pay/init/common"
	"time"
)

type ZyConfig struct {
	Zhengyuan struct {
		AppId     string `yaml:"app_id"`
		AppSecret string `yaml:"app_secret"`
		ApiUrl    string `yaml:"api_url"`
		LogPath   string `yaml:"log_path"`
	}
}

func (c *ZyConfig) LoadYaml() {
	err := yaml.Unmarshal(common.PayConf, &c)
	if err != nil {
		common.Log.Fatalf("yaml 解析失败,err:%s", err)
	}
	// 创建日志文件目录
	err = os.MkdirAll(c.Zhengyuan.LogPath, os.ModePerm)
	if err != nil {
		common.Log.Fatalf("Failed to create log directory:%s", err)
	}
	// 设置日志文件名
	logFileName := filepath.Join(c.Zhengyuan.LogPath, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		common.Log.Fatalf("Failed to open log file:%s", err)
	}

	common.Log.SetOutput(logFile)
}
