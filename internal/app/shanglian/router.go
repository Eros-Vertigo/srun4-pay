package shanglian

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiConfig "github.com/srun-soft/api-sdk/configs"
	"github.com/srun-soft/api-sdk/sdk"
	"github.com/srun-soft/pay/configs"
	"github.com/srun-soft/pay/internal/app/shanglian/config"
	"github.com/srun-soft/pay/internal/app/shanglian/middleware"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	config.C = &config.Config{}
	config.C.Load()

	apiConfig.Config = &apiConfig.APIConfig{
		Scheme:      fmt.Sprintf("%s://", config.C.Scheme),
		InterfaceIP: fmt.Sprintf("%s:8001", config.C.InterfaceIP),
		AppId:       config.C.AppId,
		AppSecret:   config.C.AppSecret,
	}
	config.API = &sdk.APIClient{}

	if *configs.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	route := gin.Default()
	//
	v1 := route.Group("/v1")
	{
		v1.POST("/check", middleware.VerifySignature(), Check)
		v1.POST("/notify", middleware.VerifySignature(), Notify)
	}

	// 启动服务
	go func() {
		err := route.Run(fmt.Sprintf(":%s", config.C.Port))
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	// 优雅的关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
