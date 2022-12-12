/*
* @Time ： 2022-12-12 22:07
* @Auth ： 张齐林
* @File ：Request_Get.go
* @IDE ：GoLand
 */
package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r:= gin.Default()
	r.GET("/get",getMsg)
	// err := r.Run("127.0.0.1:9090")
	err := r.Run(":9090")
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func getMsg(c *gin.Context) {
	name := c.Query("name")
	// 返回 String 结果容
	// c.String(http.StatusOK,"欢迎您: %s",name)
	// 返回 JSON 结果
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"data":"欢迎您: "+ name,
		"msg":"返回信息",
	})
}