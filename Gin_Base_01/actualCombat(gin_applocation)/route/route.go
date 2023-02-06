/*
* @Time ： 2023-02-06 18:03
* @Auth ： 张齐林
* @File ：route.go
* @IDE ：GoLand
 */
package route

import (
	"gin_applocation/middleware"

	"github.com/gin-gonic/gin"
)

func CpllectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoverMiddleware())
	r.POST("/api/auth/register")                         // 注册
	r.POST("/api/auth/login")                            // 登录
	r.GET("/api/auth/info", middleware.AuthMiddleware()) // 再传递数据
	return r
}
