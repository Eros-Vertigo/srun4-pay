package common

import (
	"flag"
	"fmt"
	"github.com/Eros-Vertigo/srun4-config/config"
	format "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
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
		fmt.Println("Failed to load common conf:", err)
		os.Exit(1)
	}
	// 加载yaml文件
	if *Mode == "prod" {
		PayYaml = "/srun3/etc/srun4-pay/pay.yaml"
	} else {
		PayYaml = "configs/pay.yaml"
	}
	PayConf, err = os.ReadFile(PayYaml)
	if err != nil {
		fmt.Printf("Read config [%s] failed::%s", PayYaml, err)
		os.Exit(1)
	}

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
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
		FieldsOrder:     []string{"component", "category"},
	})
}
