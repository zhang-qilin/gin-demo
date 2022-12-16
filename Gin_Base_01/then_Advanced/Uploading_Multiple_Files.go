/*
* @Time ： 2022-12-16 23:22
* @Auth ： 张齐林
* @File ：Uploading_Multiple_Files.go
* @IDE ：GoLand
 */
package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	// c.Request.ParseMultipartForm()并不能限制上传文件的大小，只是限制了上传的文件读取到内存部分的大小，如果超过了就存入了系统的临时文件中。
	// 如果需要限制文件大小，需要使用github.com/gin-contrib/size中间件，如demo中使用r.Use(limits.RequestSizeLimiter(4 << 20))限制最大4Mb。
	r.MaxMultipartMemory= 8 << 20   // 8 MB
	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()  //  获取form
		if err != nil {
			c.String(http.StatusBadRequest,"文件上传错误：%s",err)
		}
		files := form.File["file_key"]  // 上传的所有文件
		
		// 定义文件上传后存放的位置
		dst :="D:/种子/github.com/gin-demo/"
		
		// 遍历所有文件
		for _, file :=range files{
			c.SaveUploadedFile(file,dst + file.Filename)
		}
		c.String(http.StatusOK,"%d 个文件上传成功...",len(files))
	})
	
	err := r.Run()
	if err != nil {
		log.Println(err)
		return
	}
}