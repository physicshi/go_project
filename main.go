package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/zserge/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

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

// 获取在各个局域网的ip地址，转为json，作为api返回
func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs() // 获取所有网络接口
	var result []string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

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

func QrcodesController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", png)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

// 类似于TextsController的逻辑
// 获取go执行文件所在目录，创建uploads目录，拼接uploads目录+随机文件名，写入文件，返回文件路径
func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw") // 读取用户上传的文件
	if err != nil {
		log.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe) // 获取当前执行文件的
	if err != nil {
		log.Fatal(err)
	}
	filename := uuid.New().String()
	uploads := filepath.Join(dir, "uploads")
	err = os.MkdirAll(uploads, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	fullpath := path.Join("uploads", filename+filepath.Ext(file.Filename))
	fileErr := c.SaveUploadedFile(file, filepath.Join(dir, fullpath))
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
}

func main() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.POST("/api/v1/files", FilesController)
		router.GET("/api/v1/qrcodes", QrcodesController)
		router.GET("/uploads/:path", UploadsController)
		router.GET("/api/v1/addresses", AddressesController)
		router.POST("/api/v1/texts", TextsController)
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
