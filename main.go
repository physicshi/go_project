package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()
				start, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, start.Size(), "text/html", reader, nil)
			} else {
				c.Status(http.StatusNotFound)
			}

		})
		router.Run(":8080")
	}()

	var ui lorca.UI
	ui, _ = lorca.New("http://127.0.0.1:8080/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")
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
