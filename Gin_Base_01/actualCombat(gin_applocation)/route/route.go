/*
* @Time ： 2023-02-06 18:03
* @Auth ： 张齐林
* @File ：route.go
* @IDE ：GoLand
 */
package route

import (
	"gin_applocation/controller"
	"gin_applocation/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoverMiddleware())
	r.POST("/api/auth/register", controller.Register)                     // 注册
	r.POST("/api/auth/login", controller.Login)                           // 登录
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info) // 再传递数据
	return r
}
