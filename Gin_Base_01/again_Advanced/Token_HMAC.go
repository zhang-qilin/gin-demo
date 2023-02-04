/*
* @Time ： 2023-02-04 10:22
* @Auth ： 张齐林
* @File ：Token_HMAC.go
* @IDE ：GoLand
 */

// Then HMAC sinning method(HS256,HS384,HS512)  // hash 消息认证码

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Token的签名(证书签名密钥)【该密钥非常重要，如果Client端有该密钥的话，就可以自己分发密钥了】
var jwtKey = []byte("a_secret_key")

// HmacUser 客户端请求的数据
type HmacUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type MyClaims struct {
	UserId string
	// jwt.StandardClaims
	jwt.RegisteredClaims // 注意这是Google-jwt的v4版本新增的，原先(https://github.com/dgrijalva/jwt-go)是jwt.StandardClaims
}

func main() {

	r := gin.Default()
	// Token的分发
	r.POST("/getToken1", func(c *gin.Context) {
		var u HmacUser
		c.Bind(&u)
		token, err := hmacReleaseToken(u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			fmt.Println("Token生成失败...")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "token分发成功...",
			"data": token,
		})
	})
	// Token的认证
	r.POST("/checkToken1", hmacAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, "验证成功")
	})
	r.Run(":9090")

}

func hmacAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置Token的前缀(用于判断前端传递的头部是从哪个位置开始就是我们的Token)
		auth := "zhangqilin"
		// 获取authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "前缀错误...",
			})
			c.Abort()
			return
		}
		// 获取正在的Token值
		index := strings.Index(tokenString, auth+":") // 找到Token前缀对应的位置
		// 提取到真正的Token的值
		tokenString = tokenString[index+len(auth)+1:] // 真正的Token 的开始位置为：索引开始的位置+关键字的的长度(:的长度为1)
		// 对Token进行验证
		token, claims, err := hamcParseToken(tokenString)
		if err != nil || !token.Valid { // 解析错误或者过期登
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "证书无效...",
			})
			c.Abort()
			return
		}
		// claims的id和客户端提交的id是否一致
		var u HmacUser
		c.Bind(&u)
		if u.Id != claims.UserId {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 解析Token
func hamcParseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

// 分发Token
func hmacReleaseToken(u HmacUser) (string, error) {
	// 获取当前的时间，来处理Token的时效性
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 截止时间：从当前时刻算起，7天
	claims := &MyClaims{
		UserId: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "zhangqilin", // Token发布者
			Subject: "user token", // 设置主题
			// ID:       "",                               // Token的ID
			// Audience:  nil,
			// ExpiresAt: expirationTime.Unix(),           // 设置过期时间  【使用github.com/dgrijalva/jwt-go的jwt-go直接使用这种方法即可】
			// IssuedAt: time.Now().Unix(),                // 发布时间  【使用github.com/dgrijalva/jwt-go的jwt-go直接使用这种方法即可】
			// NotBefore: nil,                             // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 发布时间
			NotBefore: jwt.NewNumericDate(time.Now()),     // 生效时间 【使用Google-jwt的v4版本直接使用这种方法即可】
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 设置过期时间【使用Google-jwt的v4版本直接使用这种方法即可】
		},
	}
	// 生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 对Token加密(签名)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
