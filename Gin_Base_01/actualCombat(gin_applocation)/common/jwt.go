/*
* @Time ： 2023-02-06 18:42
* @Auth ： 张齐林
* @File ：jwt.go
* @IDE ：GoLand
 */
package common

import (
	"gin_applocation/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("a_secret_key")

type Claims struct {
	UserId string
	// jwt.StandardClaims
	jwt.RegisteredClaims // 注意这是Google-jwt的v4版本新增的，原先(https://github.com/dgrijalva/jwt-go)是jwt.StandardClaims
}

// ReleaseToken 分发证书
func ReleaseToken(u model.User) (string, error) {
	// 获取当前的时间，来处理Token的时效性
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 截止时间：从当前时刻算起，7天
	claims := &Claims{
		UserId: u.Telephone,
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

// ParseToken 解析证书
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
