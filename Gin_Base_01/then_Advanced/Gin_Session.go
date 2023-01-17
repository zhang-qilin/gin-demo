/*
* @Time ： 2023-01-17 21:23
* @Auth ： 张齐林
* @File ：Gin_Session.go
* @IDE ：GoLand
 */
package main

import (
	"net/http"
	
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	sessionName string
	sessionValue string
)

type MyOption struct {
	sessions.Options
}

func main() {

	r:= gin.Default()
	
	// 路由上加入Session中间件
	store := cookie.NewStore([]byte("session_secret"))
	r.Use(sessions.Sessions("mysession", store))
	
	r.GET("/session", func(c *gin.Context) {
		name := c.Query("name")
		if len(name) <= 0 {
			c.JSON(http.StatusBadRequest,"数据错误...")
			return
		}
		sessionName = "session_" + name
		sessionValue = "session_value_" + name
		session := sessions.Default(c)  // 获取session
		sessionData := session.Get(sessionName)
		if sessionData != sessionValue {
			// 保存session值
			session.Set(sessionName,sessionValue)
			o := MyOption{}
			o.Path = "/"
			o.MaxAge = 10   // 有效期，单位 s
			session.Options(o.Options)
			session.Save()  // 保存session
			c.JSON(http.StatusOK,"首次访问... session 已保存...")
			return
		}
		c.JSON(http.StatusOK, "访问成功... 您的session: " + sessionData.(string))
		
	})
	r.Run(":9090")

	
}

// todo:文档地址：https://github.com/gin-contrib/sessions