package common

import (
	"flag"
	format "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/srun-soft/config/config"
	"os"
	"time"
)

var (
	Log     *logrus.Logger
	Conf    *config.Config
	PayYaml string
	PayConf []byte
	err     error
	Mode    = flag.String("mode", "prod", "Mode: prod or dev")
)

func init() {
	flag.Parse()

	initLog()
	// 加载 srun|system conf
	Conf, err = config.GetConfig("conf", *Mode)
	if err != nil {
		Log.Fatalf("Failed to load common conf:%s", err)
		os.Exit(1)
	}

	Log.WithFields(logrus.Fields{
		"srun.conf":   "加载成功",
		"system.conf": "加载成功",
	}).Info()

	// 加载yaml文件
	if *Mode == "prod" {
		PayYaml = "/srun3/etc/srun4-pay/pay.yaml"
	} else {
		PayYaml = "configs/pay.yaml"
	}
	PayConf, err = os.ReadFile(PayYaml)
	if err != nil {
		Log.Fatalf("Read Config [%s] failed:%s", PayYaml, err)
		os.Exit(1)
	}
	Log.WithField("Pay配置", "加载成功")
}

// 初始化日志
func initLog() {
	Log = logrus.New()

	if *Mode == "prod" {
		Log.SetLevel(logrus.ErrorLevel)
	} else {
		Log.SetLevel(logrus.DebugLevel)
	}

	Log.SetFormatter(&format.Formatter{
		HideKeys:        false,
		TimestampFormat: time.RFC3339,
		FieldsOrder:     []string{"component", "category"},
	})

	Log.WithField("Log组件", "加载成功").Info()
}
