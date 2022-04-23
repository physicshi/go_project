package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TextsController(c *gin.Context) {
	var json struct {
		Raw string
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() // 获取当前执行文件的路径
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe) // 获取当前执行文件的目录
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()          // 生成随机文件名
		uploads := filepath.Join(dir, "uploads") // 拼接 uploads 绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  // 创建目录
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")                            // 拼接文件绝对路径
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) // 写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) // 返回文件路径
	}
}
