/*
* @Time ： 2023-02-05 9:25
* @Auth ： 张齐林
* @File ：Token_RSA.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RSA 签名实现Token

var (
	resPrivateKey  []byte // 用于存放RSA私钥
	resPublicKey   []byte // 用于存放RSA公钥
	err2_1, err2_2 error  // 用于存放读取公钥私钥发生的Error
)

// RsaUser 客户端请求的数据
type RsaUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

// RsaClaims 定义Token的结构
type RsaClaims struct {
	UserId string `json:"user_id"`
	// jwt.StandardClaims
	jwt.RegisteredClaims // 注意这是Google-jwt的v4版本新增的，原先(https://github.com/dgrijalva/jwt-go)是jwt.StandardClaims
}

func init() {
	// 读取RSA公钥密钥
	resPrivateKey, err2_1 = os.ReadFile("./rsa_key/private.pem")
	resPublicKey, err2_2 = os.ReadFile("./rsa_key/public.pem")
	if err2_1 != nil || err2_2 != nil {
		panic(fmt.Sprintf("打开密钥文件错误, err: %s %s\n", err2_1, err2_2))
		return
	}
}

// 生成Token
func rsaJwtTokenGen(id string) (any, error) {
	PrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(resPrivateKey) // 通过RSA私钥签名
	if err != nil {
		return nil, err
	}
	//
	claims := &RsaClaims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "张齐林",        // 发布者
			Subject: "user token", // 主题
			// Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 设置过期时间，截止时间：从当前时间计算：7天
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 发布时间
			// ID:        "",
		},
	}
	// 生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// 对Token进行加密(签名)
	signedString, err := token.SignedString(PrivateKey)
	return signedString, err
}

// 颁发Token
func rsaReleaseToken(u RsaUser) (any, error) {
	// 获取Token
	tokenGen, err := rsaJwtTokenGen(u.Id)
	return tokenGen, err
}

func getToken2(c *gin.Context) {
	u := RsaUser{}
	err := c.Bind(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	token, err := rsaReleaseToken(u)
	if err != nil {
		c.JSON(http.StatusOK, "生成Token失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "Token颁发成功", "data": token})

}

// 解析Token
func rsaJwtTokenRead(tokenString string) (any, error) {
	// 先获取私钥
	pem, err := jwt.ParseRSAPublicKeyFromPEM(resPublicKey)
	if err != nil {
		return nil, err
	}
	// 解析Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("解析的方法错误")
		}
		return pem, err
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && token != nil {
		return claims, nil
	}
	return nil, err

}

// Token 认证中间件(权限控制)
func rsaTokenMiddle(c *gin.Context) {
	// 设置Token的前缀(用于判断前端传递的头部是从哪个位置开始就是我们的Token)
	auth := "zhangqilin"
	// 获取authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusUnauthorized, "msg": "无效的token"})
		c.Abort()
		return
	}
	// 获取Token的下标
	index := strings.Index(tokenString, auth+":") // 找到Token前缀对应的位置
	// 提取到真正的Token的值
	tokenString = tokenString[index+len(auth)+1:] // 真正的Token 的开始位置为：索引开始的位置+关键字的的长度(:的长度为1)
	// 对Token进行验证
	claims, err := rsaJwtTokenRead(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "证书无效"})
		c.Abort()
		return
	}
	claimsValue := claims.(jwt.MapClaims) // 断言
	if claimsValue["user_id"] == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
		c.Abort()
		return
	}
	u := RsaUser{}
	err = c.Bind(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	id := claimsValue["user_id"].(string)
	if u.Id != id {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
		c.Abort()
		return
	}
	c.Next()
}

func checkToken2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "验证通过"})
}

func main() {

	r := gin.Default()
	r.POST("/getToken2", getToken2)
	r.POST("/checkToken2", rsaTokenMiddle, checkToken2)
	r.Run(":9090")

}

// TODO: RSA密钥生成工具：http://www.metools.info/code/c80.html
