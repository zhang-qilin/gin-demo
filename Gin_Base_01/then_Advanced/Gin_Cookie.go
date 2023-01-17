/*
* @Time ： 2023-01-17 20:45
* @Auth ： 张齐林
* @File ：Gin_Cookie.go
* @IDE ：GoLand
 */
package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

var (
	cookieName string
	cookieValue string
)


func CookieAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		val,_ := context.Cookie(cookieName)
		if val == ""{         // 如果是首次登录，分发cookie；浏览器保存cookie，下次登录有效
		context.SetCookie(cookieName,cookieValue,3600,"/","localhost",true,true)
		fmt.Println("Cookie已经保存完成...")
		}
	}
}

func main() {

	r:=gin.Default()
	
	// Cookie 中间件
	r.Use(CookieAuth())
	
	r.GET("/cookie", func(c *gin.Context) {
		name:=c.Query("name")
		if len(name) <= 0 {
			c.JSON(http.StatusBadRequest,"数据错误")
			return
		}
		cookieName = "cookie_" + name  // cookie的名称
		cookieValue = hex.EncodeToString([]byte(cookieName+"value"))  // cookie的值
		val, _ := c.Cookie(cookieName)
		if val == "" {
			c.String(http.StatusOK,"Cookie: %s 已经下发，下次登录有效",cookieName)
			return
		}
		c.String(http.StatusOK,"验证成功... cookie值为: %s",val)
	})
	r.Run(":9090")
	
}
