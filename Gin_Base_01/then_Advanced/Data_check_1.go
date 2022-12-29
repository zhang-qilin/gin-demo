/*
* @Time ： 2022-12-29 22:03
* @Auth ： 张齐林
* @File ：Data_check_1.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	"unicode/utf8"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)



type UserInfo struct {
	Id string `json:"id" validate:"uuid"`            // UUID  类型
	Name string `json:"name"  validate:"checkName"`  // 自定义校验
	Age uint8 `json:"age" validate:"min=0,max=130"`  // 0 <= Age <= 130
}

var validate *validator.Validate

func init()  {
	validate = validator.New()
	// 注册自定义的校验规则
	validate.RegisterValidation("checkName",checkNameFunc)
}

func checkNameFunc(f validator.FieldLevel) bool {
	count := utf8.RuneCountInString(f.Field().String())
	if count >= 2 && count <= 12 {
		return true
	}
	return false
}

func main() {
	
	u2 := uuid.NewV4()
	fmt.Printf("UUID: %s\n\n", u2)
	r:= gin.Default()
	var user UserInfo
	r.POST("/validate", func(c *gin.Context) {
		err := c.Bind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest,"请求参数错误...")
			return
		}
		// 校验
		err = validate.Struct(user)
		if err != nil {
			// 输出错误的校验值
			for _,e := range err.(validator.ValidationErrors){
				fmt.Println("错误的字段: ",e.Field())
				fmt.Println("错误的内容: ",e.Error())
				fmt.Println("错误的参数: ",e.Param())
				fmt.Println("错误具体值: ",e.Value())
				fmt.Println("错误的标签: ",e.Tag())
			}
			c.JSON(http.StatusBadRequest,"数据校验失败")
			return
		}
		c.JSON(http.StatusOK,"数据成功...")
	})
	
	r.Run()
}

// UUID 引用: github.com/satori/go.uuid
// 添加方法: go get github.com/satori/go.uuid