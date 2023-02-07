/*
* @Time ： 2023-02-06 18:08
* @Auth ： 张齐林
* @File ：RecoveryMiddleware.go
* @IDE ：GoLand
 */

package middleware

// 处理错误中间件

import (
	"fmt"
	"gin_applocation/response"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// 主要是用于获取错误
			if err := recover(); err != nil {
				response.Fail(c, nil, fmt.Sprint(err))
				c.Abort()
				return
			}
		}()
	}
}
