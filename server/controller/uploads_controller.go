package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func getUploadsDir() (uploads string) {
	exe, err := os.Executable() // 获取当前执行文件的路径
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe) // 获取当前执行文件的目录
	return filepath.Join(dir, "uploads")
}

// 手机下载文件，网络路径转为本地路径，读取本地文件，写到http响应中
func UploadsController(c *gin.Context) {
	if path := c.Param("path"); path != "" {
		target := filepath.Join(getUploadsDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}
