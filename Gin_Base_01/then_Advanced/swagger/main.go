/*
* @Time ： 2023-01-02 17:13
* @Auth ： 张齐林
* @File ：main.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	
	_ "swagger/docs"
	
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data string `json:"data"`
}

// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
func main() {
	r := gin.Default()
	// Swagger 中间件主要的作用是: 方便前端对接口进行调试。不影响接口的实际功能
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))  // 使用 Swagger 中间件
	r.GET("/login",login)
	r.POST("register",register)
	r.Run()
	
}

// @Tags 登录接口
// @Summary 登录
// @Description login
// @Accept json
// @Produce json
// @Param username query string true "用户名"
// @Param password query string false "密码"
// @Success 200 {string} json "{"code":200,"data":"{"name":"username","password":"password"}","msg":"OK"}"
// @Router /login [get]
func login(c *gin.Context) {
	userName :=c.Query("name")
	pwd := c.Query("pwd")
	fmt.Println(userName,pwd)
	res:=Response{
		Code: http.StatusOK,
		Msg:  "登录成功...",
		Data: "OK",
	}
	c.JSON(http.StatusOK,res)
}

// @Tags 注册接口
// @Summary 注册
// @Description register
// @Accept json
// @Produce json
// @Param user_name formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {string} json "{"code":200,"data":"{"name":"username","password":"password"}","msg":"OK"}"
// @Router /register [post]
func register(c *gin.Context) {
	var user User
	// err := c.Bind(&user)
	err := c.BindQuery(&user)
	if err != nil {
		fmt.Println("绑定错误: ",err)
		c.JSON(http.StatusBadRequest,"数据错误...")
		return
	}
	res:= Response{
		Code: http.StatusOK,
		Msg:  "注册成功...",
		Data: "OK",
	}
	c.JSON(http.StatusOK,res)
}
