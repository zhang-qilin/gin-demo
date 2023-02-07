/*
* @Time ： 2023-02-06 19:41
* @Auth ： 张齐林
* @File ：UserController.go
* @IDE ：GoLand
 */
package controller

import (
	"fmt"
	"gin_applocation/common"
	"gin_applocation/model"
	"gin_applocation/response"
	"gin_applocation/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 实现核心功能

// Register 注册
func Register(ctx *gin.Context) {
	var requestUser model.User
	ctx.Bind(&requestUser)
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	// 数据验证
	if len(telephone) != 11 {
		// 422 Unprocessable Entity 无法处理的请求实体
		response.Response(ctx, http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, nil, "手机号码必须为11位数")
		fmt.Println(telephone, len(telephone))
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, nil, "密码不能少于6位数")
		return
	}
	// 如果名称没有传递，给一个10位数的随机字符
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	// 判断手机号码是否存在
	if isTelephoneExist(common.DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, nil, "用户已经存在")
	}
	// 创建用户
	// 返回密码的hash值(对用户密码进行二次处理防止系统管理人员利用)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, http.StatusInternalServerError, nil, "密码加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashPassword),
	}
	common.DB.Create(&newUser) // 新增记录
	// 颁发Token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "系统异常"})
		return
	}
	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	var requestUser model.User
	ctx.Bind(&requestUser)
	// name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	// 数据验证
	if len(telephone) != 11 {
		// 422 Unprocessable Entity 无法处理的请求实体
		response.Response(ctx, http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, nil, "手机号码必须为11位数")
		fmt.Println(telephone, len(telephone))
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, nil, "密码不能少于6位数")
		return
	}
	// 依据手机号码，查询用户注册的数据记录
	var user model.User
	common.DB.Where("telephone = ? ", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": http.StatusUnprocessableEntity, "msg": "用户不存在"})
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "密码错误"})
		return
	}
	// 分发Token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "系统异常"})
		return
	}
	// 放回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{
		"user": response.ToUserDto(user.(model.User))}, "响应成功")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ? ", telephone).First(&user)
	// 如果没有查询到数据，对应uint数据，默认值为：0
	if user.ID != 0 {
		return true
	}
	return false
}
