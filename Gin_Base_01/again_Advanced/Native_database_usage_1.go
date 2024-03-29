/*
* @Time ： 2023-02-01 15:31
* @Auth ： 张齐林
* @File ：Native_database_usage_1.go
* @IDE ：GoLand
 */
package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var sqlDb *sql.DB
var sqlResponse SqlResponse

func init() {
	// 1、打开数据库
	// parseTime: 时间格式转换（查询结果为时间时，是否自动解析为时间）；
	// Loc=Local: MySQL的时区设置
	sqlStr := "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	sqlDb, err = sql.Open("mysql", sqlStr)
	if err != nil {
		fmt.Println("打开数据库出现问题:", err)
		return
	}
	// 2 尝试与数据库建立的连接(校验连接是否正确)
	err = sqlDb.Ping()
	if err != nil {
		fmt.Println("连接数据库出现问题:", err)
		return
	}
}

// SqlUser 定义结构体、客户端提交的数据
type SqlUser struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

// SqlResponse 响应体(响应Client的请求)
type SqlResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func main() {

	r := gin.Default()
	// 数据库的CRUD: 可以的对应Gin框架中的post、get、put、delete方法
	r.POST("sql/insert", insertData)
	r.GET("sql/get", getData)
	r.Run(":9090")

}

func getData(c *gin.Context) {
	name := c.Query("name")
	sqlStr := "SELECT age, address FROM user WHERE name = ?"
	var u SqlUser
	err := sqlDb.QueryRow(sqlStr, name).Scan(&u.Age, &u.Address)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "查询失败..."
		sqlResponse.Data = "ERROR!!!"
		c.JSON(http.StatusOK, sqlResponse)
		c.Abort()
		return
	}
	u.Name = name
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "查询成功..."
	sqlResponse.Data = u
	c.JSON(http.StatusOK, sqlResponse)
}

func insertData(c *gin.Context) {
	var u SqlUser
	// ===== 如果有需要的话可以添加参数校验 =====
	err := c.Bind(&u)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "参数错误..."
		sqlResponse.Data = "ERROR!!!"
		c.JSON(http.StatusOK, sqlResponse)
		c.Abort()
		return
	}
	sqlStr := "INSERT INTO user(name,age,address) VALUES(?,?,?) "
	ret, err := sqlDb.Exec(sqlStr, u.Name, u.Age, u.Address)
	if err != nil {
		fmt.Printf("INSERT FAILED, ERROR: %v/n\n", err)
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "写入数据库失败..."
		sqlResponse.Data = "ERROR!!!"
		c.JSON(http.StatusOK, sqlResponse)
		c.Abort()
		return
	}
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "写入数据库成功..."
	sqlResponse.Data = "OK!!!"
	c.JSON(http.StatusOK, sqlResponse)
	fmt.Println(ret.LastInsertId()) // 打印输出结果
}

// TODO: github.com/go-sql-driver/mysql
