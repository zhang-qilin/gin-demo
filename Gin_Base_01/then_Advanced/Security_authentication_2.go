/*
* @Time ： 2023-01-31 14:56
* @Auth ： 张齐林
* @File ：Security_authentication_2.go
* @IDE ：GoLand
 */
package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pjebs/restgate"
)

func main() {

	r := gin.Default()
	r.Use(authMiddleware2())
	r.GET("/auth2", func(c *gin.Context) {
		resData := struct {
			Code int    `json:"code"`
			Data string `json:"data"`
			Msg  string `json:"msg"`
		}{http.StatusOK, "OK", "验证通过..."}
		c.JSON(http.StatusOK, resData)
	})
	r.Run(":9090")

}

var db *sql.DB

func init() {
	db, _ = SqlDB()
}

func SqlDB() (*sql.DB, error) {
	// =================数据库连接方式1=================
	DB_TYPE := "mysql"
	DB_HOST := "127.0.0.1"
	DB_PORT := "3306"
	DB_USER := "root"
	DB_NAME := "api_secure"
	DB_PASSWORD := "123456"
	openString := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open(DB_TYPE, openString)
	return db, err

	// =================数据库连接方式2=================
	//
	// // parseTime: 时间格式转换
	// // loc=local 解决数据库时间少8小时问题
	// var err error
	// db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/api_secure?charset=utf8&parseTime=true&loc=local")
	// if err != nil {
	// 	log.Fatal("数据库打开出现了问题：", err)
	// 	return db, err
	// }
	// // 尝试与数据库建立连接（校验DNS是否正确）
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal("数据库打开出现了问题：", err)
	// 	return db, err
	//
	// }
	//
}

func authMiddleware2() gin.HandlerFunc {
	return func(context *gin.Context) {
		gate := restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Database, restgate.Config{
			DB:                 db,
			TableName:          "users",
			Key:                []string{"keys"},
			Secret:             []string{"secrets"},
			HTTPSProtectionOff: true,
		})
		nextCalled := false
		nextAdapter := func(http.ResponseWriter, *http.Request) {
			nextCalled = true
			context.Next()
		}
		gate.ServeHTTP(context.Writer, context.Request, nextAdapter)
		if nextCalled == false {
			context.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
