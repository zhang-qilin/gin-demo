/*
* @Time ： 2023-01-30 17:02
* @Auth ： 张齐林
* @File ：Security_authentication_1.go
* @IDE ：GoLand
 */
package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/pjebs/restgate"
)

func main() {
	
	r:=gin.Default()
	r.Use(authMiddleware())
	r.GET("/auth1", func(c *gin.Context) {
	resData := struct {
		Code int `json:"code"`
		Msg any `json:"msg"`
		Data any `json:"data"`
	}{http.StatusOK,"校验通过...","OK"}
	c.JSON(http.StatusOK,resData)
	})
	r.Run(":9090")
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		gate := restgate.New("X-Auth-Key",
			"X-Auth-Secret",
			restgate.Static,
			restgate.Config{
				Key:                []string{"admin", "gin"},
				Secret:             []string{"123456", "gin_ok"},
				HTTPSProtectionOff: true, // 关闭https验证
			}) // 头部key标题值，头部密钥标题值，权限源，配置项
			nextCalled := false
			nextAdapter := func (http.ResponseWriter,*http.Request){
				nextCalled = true
				c.Next()
			}
			gate.ServeHTTP(c.Writer, c.Request, nextAdapter)
			if nextCalled == false {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			
	}
}

// TODO: 文档 https://github.com/pjebs/restgate