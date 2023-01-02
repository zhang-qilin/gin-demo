/*
* @Time ： 2023-01-02 15:54
* @Auth ： 张齐林
* @File ：Data_check_2.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ValUser struct {
	Name string `json:"name" validate:"required"`       // 非空
	Age uint8 `json:"age" validate:"gte=0,lte=130"`    // 0 <= Age <= 130
	Email string `json:"email" validate:"required,email"`    // 非空、email 格式
	Address []ValAddress `json:"address" validate:"dive"` // 可以拥有多个地址
}

type ValAddress struct {
	Province string `json:"province" validate:"required"`       // 非空
	City string `json:"city" validate:"required"`           // 非空
	Phone string `json:"phone" validate:"numeric,len=11"`    // 数字类型，长度为11
}

var validate1 *validator.Validate

func init()  {
	validate1 = validator.New()  // 初始化（赋值）
}

func validateUser(u ValUser) bool {
	err := validate1.Struct(u)
	if err != nil {
		// 断言为: validator.ValidationErrors, 类型为: []FieldError
		for _,e := range err.(validator.ValidationErrors){
			fmt.Println("错误的字段",e.Field())
			fmt.Println("错误的值",e.Value())
			fmt.Println("错误的tag",e.Tag())
		}
		return false
	}
	return true
}

func main() {
	
	r := gin.Default()
	var user ValUser
	r.POST("/validate1", func(c *gin.Context) {
		
		// 测试数据的产生
		// testData(c)
		/*
		
		{
		    "name": "张齐林",
		    "age": 18,
		    "email": "3323816129@qq.com",
		    "address": [
		        {
		            "province": "广东省",
		            "city": "揭阳市",
		            "phone": "19838888930"
		        },
		        {
		            "province": "广东省",
		            "city": "揭阳市",
		            "phone": "19838888930"
		        }
		    ]
		}
		
		*/
		
		err := c.Bind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest,"参数错误，绑定失败...")
			return 
		}
		// 执行参数的校验
		if validateUser(user){
			c.JSON(http.StatusOK,"数据校验成功...")
			return
		}
		c.JSON(http.StatusBadRequest,"数据校验失败...")
	})
	
	
	r.Run()
	
}

func testData(c *gin.Context) {
	address := ValAddress{
		Province: "广东省",
		City:     "揭阳市",
		Phone:    "19838888930",
	}
	user := ValUser{
		Name:    "张齐林",
		Age:     18,
		Email:   "3323816129@qq.com",
		Address: []ValAddress{address,address},
	}
	c.JSON(http.StatusOK,user)
}


