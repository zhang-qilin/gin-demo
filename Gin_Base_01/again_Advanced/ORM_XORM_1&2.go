/*
* @Time ： 2023-02-01 20:39
* @Auth ： 张齐林
* @File ：ORM_XORM_1&2.go
* @IDE ：GoLand
 */
package main

import "github.com/go-xorm/xorm"

func init() {
	sqlStr := "root:123456@tcp(127.0.0.1:3306)/xorm?charset=utf8mb4&parseTime=True&loc=Local"
	xorm.NewEngine("mysql", sqlStr)
}

func main() {

}

//  TODO:  库项目地址：go get github.com/go-xorm/xorm      官方文档：https://xorm.io/zh/docs/
