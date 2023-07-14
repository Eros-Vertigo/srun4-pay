package routers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"srun4-pay/init/common"
	"srun4-pay/init/init_card_zy"
	"srun4-pay/internal/app/card_zy/handlers"
	"srun4-pay/internal/app/card_zy/middlewares"
	"syscall"
	"time"
)

func init() {
	conf := init_card_zy.Zy
	if *common.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(middlewares.Token())
	r.GET("/test", handlers.Test)
	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", conf.Zhengyuan.Port),
		Handler: r,
	}

	// 启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.Log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		common.Log.Fatal("Server Shutdown:", err)
	}
	common.Log.Println("Server exiting")
}
