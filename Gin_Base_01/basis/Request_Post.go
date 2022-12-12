/*
* @Time ： 2022-12-12 22:07
* @Auth ： 张齐林
* @File ：Request_Get.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r:= gin.Default()
	r.POST("/post",postMsg)
	err := r.Run(":9090")
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func postMsg(c *gin.Context) {
	// name := c.Query("name")  // 获 URL 中的数据
	name := c.DefaultPostForm("name","老张")
	fmt.Println(name)
	form, b := c.GetPostForm("name")
	fmt.Println(form,b)

	c.JSON(http.StatusOK,"欢迎您: "+name)
}