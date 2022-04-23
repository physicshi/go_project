package controller

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
