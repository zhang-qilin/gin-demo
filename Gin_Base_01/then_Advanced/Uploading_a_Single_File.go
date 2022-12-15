/*
* @Time ： 2022-12-16 0:12
* @Auth ： 张齐林
* @File ：Uploading_a_Single_File.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("fileName")
		if err != nil {
			c.String(http.StatusBadRequest,"文件上传错误")
		}
		// 定义服务端要把客户端上传的文件保存到上面路径下(存储路径)
		dst :="D:/种子/github.com/gin-demo/"
		// 存储文件
		c.SaveUploadedFile(file,dst+file.Filename)
		c.String(http.StatusOK,fmt.Sprintf("%s 上传成功",file.Filename))
		
	})
	err := r.Run()
	if err != nil {
		log.Fatalf("err: %s.\n", err)
		return
	}
}