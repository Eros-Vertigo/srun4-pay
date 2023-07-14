package init_card_zy

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"srun4-pay/init/common"
	"time"
)

type zyConfig struct {
	Zhengyuan struct {
		AppId     string `yaml:"app_id"`
		AppSecret string `yaml:"app_secret"`
		ApiUrl    string `yaml:"api_url"`
		LogPath   string `yaml:"log_path"`
	}
}

var Zy *zyConfig

func init() {
	err := yaml.Unmarshal(common.PayConf, &Zy)
	if err != nil {
		fmt.Println("yaml 解析失败", err)
		os.Exit(1)
	}
	// 创建日志文件目录
	err = os.MkdirAll(Zy.Zhengyuan.LogPath, os.ModePerm)
	if err != nil {
		common.Log.Fatalf("Failed to create log directory:%s", err)
	}
	// 设置日志文件名
	logFileName := filepath.Join(Zy.Zhengyuan.LogPath, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		common.Log.Fatalf("Failed to open log file:%s", err)
	}

	common.Log.SetOutput(logFile)
}
