/*
* @Time ： 2023-02-06 18:12
* @Auth ： 张齐林
* @File ：AuthMiddleware.go
* @IDE ：GoLand
 */

package middleware

import (
	"gin_applocation/common"
	"gin_applocation/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Token认证中间件(权限控制)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := "zhangqilin"
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		// 无效的Token
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "前缀错误...",
			})
			ctx.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":") // 找到Token前缀对应的位置
		tokenString = tokenString[index+len(auth)+1:] // 真正的Token 的开始位置为：索引开始的位置+关键字的的长度(:的长度为1)
		// 对Token进行验证
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { // 解析错误或者过期登
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "证书无效..."})
			ctx.Abort()
			return
		}
		// 验证通过后获取claim中的userId
		userId := claims.UserId
		// 判定
		var user model.User
		common.DB.First(&user, userId)
		if user.ID == 0 { // 如果没有读取到内容，说明Token值有误
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		ctx.Set("user", user) // 将key-value值存储到context中
		ctx.Next()
	}
}
