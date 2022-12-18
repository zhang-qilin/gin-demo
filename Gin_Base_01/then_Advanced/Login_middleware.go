/*
* @Time ： 2022-12-18 23:58
* @Auth ： 张齐林
* @File ：Login_middleware.go
* @IDE ：GoLand
 */
package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(AuthMiddleware())
	
	r.GET("/login", func(context *gin.Context) {
		// 获取用户，它是由 BasicAuth 中间件设置的
		user:= context.MustGet(gin.AuthUserKey).(string)
		context.JSON(http.StatusOK,"登录成功..." + "欢迎您: " + user)
	})
	
	err := r.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func AuthMiddleware() gin.HandlerFunc {
	// 初始化用户
	accounts := gin.Accounts{ // gin.Accounts 是 map[string]string
		"admin":"adminpw",
		"system":"systempw",
	}
	
	// 动态添加用户
	accounts["root"]="rootpw"
	
	// 将用户添加到登录中间件中
	auth := gin.BasicAuth(accounts)
	return auth
}