/*
* @Time ： 2022-12-21 23:03
* @Auth ： 张齐林
* @File ：Programs_on_multiple_servers.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	// 定义路由组
	g errgroup.Group
)

func main() {
	// 服务器1:
	server01 := &http.Server{
		Addr: ":9091",
		Handler: router01(),
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
	// 服务器2:
	server02 := &http.Server{
		Addr: ":9092",
		Handler: router02(),
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
	
	// 开启服务
	g.Go(func() error {
		return server01.ListenAndServe()
	})
	g.Go(func() error {
		return server02.ListenAndServe()
	})
	if err := g.Wait(); err!= nil{
		fmt.Printf("执行失败, err: %s\n", err)
	}
	
}

func router01() http.Handler {
	r1 :=gin.Default()
	r1.GET("/MyServer01", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"msg":"服务器程序1的相对应",
		},
		)
	})
	return r1
}
func router02() http.Handler {
	r2 :=gin.Default()
	r2.GET("/MyServer02", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"msg":"服务器程序2的相对应",
		},
		)
	})
	return r2
}
