package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

func main() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.String(200, "Hello World")
		})
		router.Run(":8080")
	}()

	var ui lorca.UI
	ui, _ = lorca.New("http://127.0.0.1:8080", "", 800, 600, "--disable-sync", "--disable-translate")
	// 处理中断、中止信号
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	// select会等待第一个可以读或者写的ch进行操作
	select {
	case <-ui.Done():
	case <-chSignal:
	}
	ui.Close()
}
