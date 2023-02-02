/*
* @Time ： 2023-02-01 20:39
* @Auth ： 张齐林
* @File ：ORM_XORM_1&2.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var x *xorm.Engine
var xormResponse XormResponse

// Stu 定义结构体(XROM支持双向映射)：没有表，会进行创建【结构体名称必须与表名称一致】
type Stu struct {
	Id      int64     `xorm:"pk autoincr" json:"id"` // 指定主键并自增
	StuNum  string    `xorm:"unique" json:"stu_num"` // 设置唯一值
	Name    string    `json:"name"`
	Age     int       `json:"age"`
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}

// XormResponse 应答体
type XormResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func init() {
	sqlStr := "root:123456@tcp(127.0.0.1:3306)/xorm?charset=utf8&parseTime=true&loc=Local"
	var err error
	// 1、 创建数据库引擎
	x, err = xorm.NewEngine("mysql", sqlStr)
	if err != nil {
		fmt.Println("数据库连接失败: ", err)
		return
	}
	// 默认xorm采用Local时区，所以默认调用的time.Now()会先被转换成对应的时区。
	// 要改变xorm的时区，则需要用到以下语句
	x.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	// 2、创建或同步表(名称为Stu)
	err = x.Sync2(new(Stu))
	if err != nil {
		fmt.Println("数据库同步失败: ", err)
	}
}

func xormDelete(c *gin.Context) {
	StuNum := c.Query("stu_num")
	// 先查询
	var stus []Stu
	err := x.Where("stu_num = ?", StuNum).Find(&stus)
	if err != nil || len(stus) <= 0 {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "查询失败或查询的数据不存在"
		xormResponse.Data = "ERROR"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	// 再删除
	affected, err := x.Where("stu_num = ?", StuNum).Delete(&Stu{})
	if err != nil || affected <= 0 {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "删除失败或删除的数据不存在"
		xormResponse.Data = "ERROR"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "删除成功"
	xormResponse.Data = "OK"
	c.JSON(http.StatusOK, xormResponse)
	return
}

func xormUpdateData(c *gin.Context) {
	var s Stu
	err := c.Bind(&s)
	if err != nil {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "参数错误"
		xormResponse.Data = "ERROR"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	// 先查找
	var stus []Stu
	err = x.Where("stu_num = ?", s.StuNum).Find(&stus)
	if err != nil || len(stus) <= 0 {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "查询失败或查询数据不存在"
		xormResponse.Data = "ERROR"
		c.JSON(http.StatusOK, xormResponse)
	}
	// 再修改
	affected, err := x.Where("stu_num = ?", s.StuNum).Update(&Stu{Name: s.Name, Age: s.Age})
	if err != nil || affected <= 0 {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "更新失败"
		xormResponse.Data = err
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "更新成功"
	xormResponse.Data = "OK"
	c.JSON(http.StatusOK, xormResponse)
}

func xormGetMulData(c *gin.Context) {
	name := c.Query("name")
	var stus []Stu
	err := x.Where("name = ?", name).And("age > 20").Limit(10, 0).Asc("age").Find(&stus)
	if err != nil {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "查询错误"
		xormResponse.Data = "ERROR"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "读取成功"
	xormResponse.Data = stus
	c.JSON(http.StatusOK, xormResponse)
}

func xormGetData(c *gin.Context) {
	StuNam := c.Query("stu_num")
	var stus []Stu
	err := x.Where("stu_num = ?", StuNam).Find(&stus)
	if err != nil {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "查询错误"
		xormResponse.Data = "ERROR"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "读取成功"
	xormResponse.Data = stus
	c.JSON(http.StatusOK, xormResponse)
}

func xormInsert(c *gin.Context) {
	var s Stu
	err := c.Bind(&s)
	if err != nil {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "参数错误"
		xormResponse.Data = "ERROR"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	affected, err := x.Insert(s)
	if err != nil || affected <= 0 {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "写入失败"
		xormResponse.Data = err
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "写入成功"
	xormResponse.Data = "OK"
	c.JSON(http.StatusOK, xormResponse)
	fmt.Println(affected)
}

func main() {
	r := gin.Default()
	// 数据库的CRUD 对应Gin的Post、Get、Put、Delete方法
	r.POST("xorm/insert", xormInsert)    // 添加数据
	r.GET("xorm/get", xormGetData)       // 查询数据(单条记录)
	r.GET("xorm/mulget", xormGetMulData) // 查询数据(多条记录)
	r.PUT("xorm/update", xormUpdateData) // 更新数据
	r.DELETE("xorm/delete", xormDelete)  // 删除数据
	r.Run(":9090")
}

//  TODO:  库项目地址：go get github.com/go-xorm/xorm      官方文档：https://xorm.io/zh/docs/
