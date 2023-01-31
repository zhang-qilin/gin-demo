/*
* @Time ： 2023-01-31 17:08
* @Auth ： 张齐林
* @File ：Log_frame_logrus_1.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 全局创建一个log实例
var log = logrus.New()

func initLorus() error {
	log.Formatter = &logrus.JSONFormatter{} // 设置为JSON格式的日志
	// log.Formatter = &logrus.TextFormatter{} // 设置为Text格式的日志
	file, err := os.OpenFile("./gin_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("创建文件/打开文件失败...")
		return err
	}
	// 设置log的默认文件输出
	log.Out = file
	// 设置log的默认终端输出
	// log.Out = os.Stdout
	// 设置gin框架的版本为: 发布版本
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Out
	// 设置日志级别
	log.Level = logrus.InfoLevel
	return nil
}

func main() {

	err := initLorus()
	if err != nil {
		fmt.Println(err)
		return
	}
	r := gin.Default()
	r.GET("/logrus", func(c *gin.Context) {
		log.WithFields(logrus.Fields{
			"url":    c.Request.RequestURI,
			"method": c.Request.Method,
			"params": c.Query("name"),
			"IP":     c.ClientIP(),
		}).Info()
		resData := struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data string `json:"data"`
		}{http.StatusOK, "响应成功...", "OK!"}
		c.JSON(http.StatusOK, resData)
	})
	r.Run(":9090")

}
