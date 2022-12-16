/*
* @Time ： 2022-12-17 0:14
* @Auth ： 张齐林
* @File ：Custom_middleware.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 如果需要用到没有任何中间件的路由引擎时可以使用以下的 gin.New() 方法
	// r:=gin.New()
	r.Use(Midleware())  // 使用User的话是属于全局中间件的定义使用，如果只需要当单独的方法使用中间件的话(也就是说局部中间件)，直接在方法后面的路由添加即可，如下
	// r.GET("/middleware",Midleware(), func(c *gin.Context) {
	r.GET("/middleware", func(c *gin.Context) {
	fmt.Println("服务端开始执行...")
		
		name := c.Query("name")
		ageStr := c.Query("age")
		age, _ := strconv.Atoi(ageStr)
		log.Println(name, age)
		
		res := struct {
			Name string `json:"name"`
			Age int `json:"age"`
		}{name,age,}
		c.JSON(http.StatusOK,res)
		
	})
	
	r.Run()
}

func Midleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("中间件开始执行...")
		name := c.Query("name")
		ageStr := c.Query("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,"用户数据的数据有错误，age 的值不是整数")
			
			// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
			// Let's say you have an authorization middleware that validates that the current request is authorized.
			// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
			// for this request are not called.
			c.Abort()
			return
		}
		if age <0 || age> 100{
			c.AbortWithStatusJSON(http.StatusBadRequest,"用户数据的数据有错误，age 的值不是整数")
			c.Abort()
			return
		}
		if len(name) <6 || len(name) >12{
			c.AbortWithStatusJSON(http.StatusBadRequest,"用户名只能是6-12位数")
			c.Abort()
			return
		}
		// 执行后续的操作
		c.Next()
		fmt.Println(name, age)
	}
}