/*
* @Time ： 2022-12-15 23:20
* @Auth ： 张齐林
* @File ：Multi_form_rendering.go
* @IDE ：GoLand
 */
package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	// JSON 格式输出
	r.GET("/json", func(c *gin.Context) {
		
		c.JSON(http.StatusOK,gin.H{
			"html":"<b>Hello Gin</b>",
		})
		
	})
	// HTML 格式输出
	r.GET("/someHTML", func(c *gin.Context) {
		c.PureJSON(http.StatusOK,gin.H{
			"html":"<b>Hello Gin</b>",
		})
	})
	
	// XML 格式输出(XML 渲染)
	r.GET("/someXML", func(c *gin.Context) {
		type Message struct {
			Name string
			Msg string
			Age int
		}
		info := Message{
			Name: "张齐林",
			Msg:  "Hello",
			Age:  88,
		}
		c.XML(http.StatusOK,info)
	})
	
	// YAML 格式输出(YAML 渲染)
	r.GET("/someYAML", func(c  *gin.Context) {
		c.YAML(http.StatusOK,gin.H{
			"message":"Gin 框架的多形式渲染",
			"status":http.StatusOK,
		})
	})
	
	r.Run(":9090")
}

