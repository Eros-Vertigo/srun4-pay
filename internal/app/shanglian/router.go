package shanglian

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	C = &Config{}
	C.Load()

	route := gin.Default()
	//
	v1 := route.Group("/v1")
	{
		v1.POST("/check", VerifySignature(), Check)
		v1.POST("/notify", VerifySignature(), Notify)
	}

	// 启动服务
	go func() {
		err := route.Run(":8890")
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	// 优雅的关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
