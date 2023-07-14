package init_card_zy

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"srun4-pay/init/common"
	"time"
)

const (
	CacheAccessToken = "srun4_pay:cache_access_token"
)

var (
	Zy *ZyConfig
)

type ZyConfig struct {
	Zhengyuan struct {
		AppId     string `yaml:"app_id"`
		AppSecret string `yaml:"app_secret"`
		ApiUrl    string `yaml:"api_url"`
		LogPath   string `yaml:"log_path"`
		Port      string `yaml:"port" default:"8890"`
	}
}

func init() {
	Zy = &ZyConfig{}
	Zy.LoadYaml()
}

// LoadYaml 加载「正元一卡通」配置
// 设置细分日志目录
func (c *ZyConfig) LoadYaml() {
	err := yaml.Unmarshal(common.PayConf, &c)
	if err != nil {
		common.Log.WithField("yaml 解析失败,err:%s", err).Fatal()
	}
	common.Log.WithField("正元一卡通配置", "加载成功").Info()
	// 创建日志文件目录
	err = os.MkdirAll(c.Zhengyuan.LogPath, os.ModePerm)
	if err != nil {
		common.Log.WithField("Failed to create log directory:%s", err).Fatal()
	}
	// 设置日志文件名
	logFileName := filepath.Join(c.Zhengyuan.LogPath, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		common.Log.WithField("Failed to open log file:%s", err).Fatal()
	}

	if *common.Mode == "prod" {
		mw := io.MultiWriter(os.Stdout, logFile)
		common.Log.SetOutput(mw)
	} else {
		common.Log.SetOutput(os.Stdout)
	}
}
