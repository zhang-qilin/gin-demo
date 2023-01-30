/*
* @Time ： 2023-01-30 15:00
* @Auth ： 张齐林
* @File ：Gin_Https.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

type HttpRes struct {
	Code int `json:"code"`
	Result string `json:"result"`
}

func main() {
	r:=gin.Default()
	r.Use(httpsHandel())   // https 对应的中间件
	r.GET("/https_test", func(c *gin.Context) {
		fmt.Println(c.Request.Host)
		c.JSON(http.StatusOK,HttpRes{
			Code:   http.StatusOK,
			Result: "测试成功...",
		})
	})
	r.RunTLS(":9090","SSL_TLS/zhangqilin.crt","SSL_TLS/zhangqilin.key")
}

func httpsHandel() gin.HandlerFunc {
	return func(context *gin.Context) {
		secureMiddle := secure.New(secure.Options{
			SSLRedirect: true,  // 只允许https请求
			// SSLHost:   // http到https的重定向
			STSSeconds: 1536000,  // Strict-Transport-Security header 的时效：1年
			// 以下部分为HTTP请求头内容，为固定写法
			STSIncludeSubdomains: true,  // IncludeSubdomains will be appended to the Strict-Transport-Security header
			STSPreload: true, // STS Preload(预加载)
			FrameDeny: true,  // X-Frame-Options 有三个值： DENY(表示该页面不允许在frame中展示，即便是在相同域名中的页面中)
			ContentTypeNosniff: true,  // 禁用浏览器的类型猜测行为防止基于 MIME 类混淆的攻击
			BrowserXssFilter: true,  // 启用XSS保护，并在检查到XSS攻击时，停止渲染页面
			// IsDevelopment: true,  // 开发模式
		})
		err := secureMiddle.Process(context.Writer, context.Request)
		// 如果不安全，终止
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest,"数据不安全")
			return
		}
		// 如果是重定向，终止
		if ststus := context.Writer.Status();ststus>300&&ststus<399{
			context.Abort()
			return
		}
		context.Next()
	}
}

// TODO: 测试证书生成工具 https://keymanager.org/
// TODO: 中间件对应的包 github.com/unrolled/secure
