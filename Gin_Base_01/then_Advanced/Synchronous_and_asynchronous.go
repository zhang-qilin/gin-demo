/*
* @Time ： 2022-12-21 22:30
* @Auth ： 张齐林
* @File ：Synchronous_and_asynchronous.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		sysnc(c)
		c.JSON(http.StatusOK,">>> 主程序(主 Go 程)同步已执行... <<<")
	})
	
	r.GET("/async", func(c *gin.Context) {
		for i := 0; i < 6; i++ {
			cCp := c.Copy()
			go async(cCp,i)
		}
		c.JSON(http.StatusOK,"~~ 主程序(主 Go 程)同步已执行... ~~")
	})
	
	err := r.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func async(cp *gin.Context, i int) {
	fmt.Println("第" + strconv.Itoa(i) + "个 Go 程开始执行..." + cp.Request.URL.Path)
	time.Sleep(time.Second * 3)
	fmt.Println("第" + strconv.Itoa(i) + "个 Go 程执行结束..." + cp.Request.URL.Path)
}

func sysnc(c *gin.Context) {
	println("执行同步任务: " + c.Request.URL.Path)
	time.Sleep(time.Second * 3)
	println("同步任务执行完成...")
}