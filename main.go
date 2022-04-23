package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/physicshi/go_project/server"
	"github.com/zserge/lorca"
)

func main() {
	go func() {
		server.RunGinServer()
	}()

	var ui lorca.UI
	ui, _ = lorca.New("http://127.0.0.1:27149/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")
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
