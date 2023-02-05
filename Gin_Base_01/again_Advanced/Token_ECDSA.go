/*
* @Time ： 2023-02-05 14:01
* @Auth ： 张齐林
* @File ：Token_ECDSA.go
* @IDE ：GoLand
 */
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RSA 非对称加密签名

var (
	err3          error
	eccPrivateKey *ecdsa.PrivateKey
	eccPublicKey  *ecdsa.PublicKey
)

// EcdsaUsre 客户端提交的数据
type EcdsaUsre struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

// EcdsaClaims Token结构体
type EcdsaClaims struct {
	UserId string `json:"user_id"`
	// jwt.StandardClaims
	jwt.RegisteredClaims // 注意这是Google-jwt的v4版本新增的，原先(https://github.com/dgrijalva/jwt-go)是jwt.StandardClaims
}

func init() {
	eccPrivateKey, eccPublicKey, err3 = getEcdsaKey(2)
	if err3 != nil {
		panic(err3)
		return
	}
}

// ecdsa密钥生成
func getEcdsaKey(keyType int) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var (
		err   error // 1：224  2：256
		prk   *ecdsa.PrivateKey
		pub   *ecdsa.PublicKey
		curve elliptic.Curve // 椭圆曲线
	)
	switch keyType {
	case 1:
		curve = elliptic.P224()
	case 2:
		curve = elliptic.P256()
	case 3:
		curve = elliptic.P384()
	case 4:
		curve = elliptic.P521()
	default:
		err = errors.New("输入的签名Key类型错误！Key取值：\n 1：椭圆曲线224 \n 2：椭圆曲线256 \n 3：椭圆曲线384 \n 4：椭圆曲线521 \n")
		return nil, nil, err
	}
	// 获取私钥
	prk, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	// 公钥赋值(通过私钥来获取公钥)
	pub = &prk.PublicKey
	return prk, pub, err

}

// 颁发Token
func ecdsaReleaseToken(u EcdsaUsre) (any, error) {
	claims := &EcdsaClaims{
		UserId: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "张齐林",                                                  // 发布者
			Subject:   "user token",                                           // 主题
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 设置过期时间，截止时间：从当前时间计算：7天
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 发布时间
		},
	}
	// 生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	// Token进行加密(签名)
	signedString, err := token.SignedString(eccPrivateKey)
	return signedString, err
}

func ecdsaJwtTokenRead(tokenString string) (any, error) {
	myToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v\n", token.Method)
		}
		return eccPublicKey, nil
	})
	if claims, ok := myToken.Claims.(jwt.MapClaims); ok && myToken.Valid {
		return claims, nil
	}
	return nil, err
}

// Token 认证中间件(权限控制)
func ecdsaTokenMiddle(c *gin.Context) {
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
	claims, err := ecdsaJwtTokenRead(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}
	claimsValue := claims.(jwt.MapClaims)
	if claimsValue["user_id"] == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
		c.Abort()
		return
	}
	u := &EcdsaUsre{}
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

func getToken3(c *gin.Context) {
	u := EcdsaUsre{}
	err := c.Bind(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	// 分发Token
	token, err := ecdsaReleaseToken(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, "生成Token错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "授权成功",
		"data": token,
	})
}

func checkToken3(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "验证通过",
	})
}

func main() {
	r := gin.Default()
	r.POST("/getToken3", getToken3)
	r.POST("/checkToken3", ecdsaTokenMiddle, checkToken3)
	r.Run(":9090")
}
