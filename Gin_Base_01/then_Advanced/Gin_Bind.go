/*
* @Time ： 2022-12-22 22:43
* @Auth ： 张齐林
* @File ：Gin_Bind.go
* @IDE ：GoLand
 */
package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type Login struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remark string `json:"remark"`
}

func main() {
	
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		
		var login Login
		err := c.Bind(&login)
		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":"绑定失败,参数错误",
				"data":err.Error(),
			})
			return
		}
		if login.UserName == "user" && login.Password == "123456"{
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":"登录成功...",
			"data":"OK",
		})
			return
		}
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":"登录失败...",
			"data":"error",
		})
	})
	
	err := r.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
