/*
* @Time ： 2023-02-05 16:34
* @Auth ： 张齐林
* @File ：main.go
* @IDE ：GoLand
 */
package main

import (
	"gin_applocation/common"
	"gin_applocation/route"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	InintConfig()
	common.InitDB()
	r := gin.Default()
	r = route.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		r.Run(":" + port)
	} else {
		r.Run() // 默认端口号:8080
	}
}

func InintConfig() {
	workDir, _ := os.Getwd()                 // 获取目录对应的路径
	viper.SetConfigName("application")       // 配置文件名
	viper.SetConfigType("yaml")              // 配置文件类型
	viper.AddConfigPath(workDir + "/config") // 执行go run 对应的路径配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
		return
	}
}
