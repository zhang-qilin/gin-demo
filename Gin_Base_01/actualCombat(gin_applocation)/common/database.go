/*
* @Time ： 2023-02-05 16:51
* @Auth ： 张齐林
* @File ：database.go
* @IDE ：GoLand
 */
package common

import (
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	hots := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	// dirverName := viper.GetString("datasource.dirverName")
	databases := viper.GetString("datasource.databases")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		hots,
		port,
		databases,
		charset,
		url.QueryEscape(loc))
	fmt.Println(args)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Failed To Connect database, err: " + err.Error())
	}
	DB = db
	return db
}
