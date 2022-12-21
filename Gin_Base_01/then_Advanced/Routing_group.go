/*
* @Time ： 2022-12-21 23:43
* @Auth ： 张齐林
* @File ：Routing_group.go
* @IDE ：GoLand
 */
package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type ResGroup struct {
	Data string
	Path string
}

func main() {
	
	router := gin.Default()
	// 路由分组1
	v1 := router.Group("/v1")  // 路由分组1(一级路径)
	{
		r := v1.Group("/user")     // 路由分组(二级路径)
		r.GET("/login",login)      // 响应请求: v1/user/login
		// 路由分组(三级路径)
		r2 := r.Group("/showInfo")  // /v1/user/showInfo
		r2.GET("/abstract",abstract)   // 响应请求: v1/user/showInfo/abstract
		r2.GET("detail",detail)   // 响应请求：v1/user/showInfo/detail
		
	}
	// 路由分组2
	v2 := router.Group("/v2")
	{
		v2.GET("/other",other)
	}
	
	err := router.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func other(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "other",
		Path: c.Request.URL.Path,
	})
}

func detail(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "detail",
		Path: c.Request.URL.Path,
	})
}

func abstract(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "abstract",
		Path: c.Request.URL.Path,
	})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "login",
		Path: c.Request.URL.Path,
	})
}
