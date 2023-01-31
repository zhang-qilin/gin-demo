/*
* @Time ： 2023-01-31 18:26
* @Auth ： 张齐林
* @File ：Log_frame_logrus_2.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logFilePath = "./"         // 日志文件保存的路径
	logFileName = "system.log" // 日志文件的名称
)

func main() {

	r := gin.Default()
	r.Use(logMiddleware())
	r.GET("/logrus2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "响应成功...",
			"data": "OK！",
		})
	})
	r.Run(":9090")
}

func logMiddleware() gin.HandlerFunc {
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	// 实例化logrus
	logger := logrus.New()
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置输出到文件中
	logger.Out = file
	// 实现rotatelogs，设置文件分割
	logsWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+"%Y%m%d.log",
		// 生成软连接，指向最新的日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour), // 以 Hour 为单位的整数
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(60*time.Second),
	)
	// hook 机制的设置
	writerMap := lfshook.WriterMap{
		logrus.DebugLevel: logsWriter,
		logrus.InfoLevel:  logsWriter,
		logrus.WarnLevel:  logsWriter,
		logrus.ErrorLevel: logsWriter,
		logrus.FatalLevel: logsWriter,
		logrus.PanicLevel: logsWriter,
	}
	// 给logrus添加hook机制
	logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",

		// 设置日志输出的格式，以下可暂时不使用
		// DisableTimestamp:  false,
		// DisableHTMLEscape: false,
		// DataKey:           "",
		// FieldMap:          nil,
		// CallerPrettyfier:  nil,
		// PrettyPrint:       false,
		//
	}))
	return func(c *gin.Context) {
		c.Next()
		// 获取相关信息
		// 请求方式(方法)
		method := c.Request.Method
		// 请求路由
		reqUrl := c.Request.URL
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"client_ip":   clientIP,
			"req_method":  method,
			"req_url":     reqUrl,
		}).Info()
	}

}
