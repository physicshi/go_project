package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/physicshi/go_project/server/controller"
	"github.com/physicshi/go_project/server/ws"
)

//go:embed frontend/dist/*
var FS embed.FS

func RunGinServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	hub := ws.NewHub()
	go hub.Run()
	router.StaticFS("/static", http.FS(staticFiles))
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.POST("/api/v1/files", controller.FilesController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.GET("/uploads/:path", controller.UploadsController)
	router.GET("/api/v1/addresses", controller.AddressesController)
	router.POST("/api/v1/texts", controller.TextsController)
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
	router.Run(":27149")
}
