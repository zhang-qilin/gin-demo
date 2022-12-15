/*
* @Time ： 2022-12-15 23:24
* @Auth ： 张齐林
* @File ：Return_third-party_data.go
* @IDE ：GoLand
 */
package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.GET("/GetOtherData", func(c *gin.Context) {
		// url:="https://www.baidu.com"
		url:="https://avatars.githubusercontent.com/u/53826118?v=4"
		response, err := http.Get(url)
		if err != nil || response.StatusCode!= http.StatusOK{
			c.Status(http.StatusServiceUnavailable)  // 应答 Client
			return
		}
		body := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		// 数据写入响应体
		c.DataFromReader(http.StatusOK,contentLength,contentType,body,nil)
	})
	
	r.Run(":9090")
}