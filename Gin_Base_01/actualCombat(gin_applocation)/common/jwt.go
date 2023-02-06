/*
* @Time ： 2023-02-06 18:42
* @Auth ： 张齐林
* @File ：jwt.go
* @IDE ：GoLand
 */
package common

import "github.com/golang-jwt/jwt/v4"

var jwtKey = []byte("a_secret_key")

type MyClaims struct {
	UserId string
	// jwt.StandardClaims
	jwt.RegisteredClaims // 注意这是Google-jwt的v4版本新增的，原先(https://github.com/dgrijalva/jwt-go)是jwt.StandardClaims
}

func ParseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
