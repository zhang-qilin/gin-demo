/*
* @Time ： 2023-02-06 15:52
* @Auth ： 张齐林
* @File ：response.go
* @IDE ：GoLand
 */
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 封装响应体
func Response(ctx *gin.Context, httpStatus, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

// Success 响应成功的响应体
func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, http.StatusOK, data, msg)
}

// Fail 响应失败的响应体
func Fail(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, http.StatusUnprocessableEntity, data, msg)
}
