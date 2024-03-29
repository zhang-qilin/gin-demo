/*
* @Time ： 2023-02-06 18:11
* @Auth ： 张齐林
* @File ：CORSMiddleware.go
* @IDE ：GoLand
 */

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 跨域中间件
// (跨域，指的是浏览器不能执行其他网站的脚本。
// 它是由浏览器的同源策略造成的，是浏览器对JavaScript施加的安全限制)
func CORSMiddleware() gin.HandlerFunc { // CORS 是跨域资源分享(Cross-Origin Resource Sharing) 中间件
	return func(ctx *gin.Context) {
		// 指定允许其他域名访问(Access-Control-Allow-Origin)
		// ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 跨域：CORS(跨来源资源共享)策略
		// 预检结果缓存时间
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 允许请求的类型(GET、POST等...)
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		// 允许的请求头字段
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		// 是否允许后续请求携带认证信息(Cookies)，该值只能是true,负责不返回
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}
