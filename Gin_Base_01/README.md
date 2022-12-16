# Gin自学笔记

官方是怎么描述Gin框架的

- Go语言最快的全功能Web框架。晶莹剔透。

- Gin是一个使用Go语言开发的Web框架。

- 它提供类似Martini的API，但性能更佳，速度提升高达40倍。

- 如果你是性能和高效的追求者, 你会爱上 Gin。

![image-20221214232322247](README.assets/image-20221214232322247.png)

## 总体学习流程

![image-20221212211824336](./README.assets/image-20221212211824336.png)

## 一、Get 请求

HTTP 是超文本传输协议。

其定义了客户端和服务端之间文本传输的规范。



### Get 方法

Get 方法：主要用于从指定资源中请求数据。

使用 Get 方法的请求应该是只是检索数据，并且不应对数据产生其它影响。

### Gin 如何实现 Http 的Get 方法

#### 准备工作

##### 1、下载并且安装 Gin

```bash
go get -u github.com/gin-gonic/gin
```

##### 2、将 Gin 导入代码中

```go
import "github.com/gin-gonic/gin"
import "net/http"
```

#### 核心代码

```go
r := gin.Default() // 获取路由引擎
r.GET()            // 执行 HTTP 协议的 Get 请求
```

##### 关于 Get 的一些问题

安全性：因 Get 请求的不安全性，在处理敏感数据时，绝不可以使用 Get 请求，因为数据在URL中对所有人都是可见的。

数据长度：当发送数据时，Get 方法向 URL 添加数据源；URL 的长度是受限制的（ URL 的最大长度是 2048 个字符）。

因此需用到下面说到的Post方法

## 二、  Post 方法

### 引入 Post 方法

1. Post 方法用于将数据发送到服务器以创建或跟新资源
2. Post 比 Get 更安全
3. 对于数据长度没有限制

在客户机和服务器之间进行请求 - 响应时，两种最常用的方法就是：`Get` 和 `Post`

Get - 从指定的资源请求数据

Post - 向指定的资源提交要被处理的数据

最直观的区别就是 Get 把参数包含在 URL 中，Post 通过request body 传递参数

### Gin 如何使用 Post 请求

#### 核心代码

```go
r := gin.Default()
r.Post()
```

## 三、重定向

将指定的网络请求重写定个方向，使其跳转到指定的其它位置(网站)

通过重定向来完成网页、网址的自动跳转

可分为：一般重定向 与 路由重定向

#### 核心代码

```go
// 重定向
// 一般重定向：重定向到外部网络
	r.GET("/redirect1", func(c *gin.Context) {
		// 重定向到本人 GitHub 首页，获取到 GitHub 本人首页的数据
		// 重定向状态码：StatusMovedPermanently
		url := "https://github.com/zhang-qilin"
		c.Redirect(http.StatusMovedPermanently,url)
	})


// Gin 路由重定向，使用如下的HandleContex
// 路由重定向：重定向到具体的路由
	r.GET("/redirect2", func(c *gin.Context) {
		c.Request.URL.Path="/TestRedirect"
		r.HandleContext(c)
	})
	// 路由：http://localhost:9090/TestRedirect
	r.GET("/TestRedirect", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"msg":"响应成功",
		})
	})
```

Demo

```go
/*
* @Time ： 2022-12-12 23:03
* @Auth ： 张齐林
* @File ：Redirection.go
* @IDE ：GoLand
 */
package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	// 一般重定向：重定向到外部网络
	r.GET("/redirect1", func(c *gin.Context) {
		// 重定向到本人 GitHub 首页，获取到 GitHub 本人首页的数据
		// 重定向状态码：StatusMovedPermanently
		url := "https://github.com/zhang-qilin"
		c.Redirect(http.StatusMovedPermanently,url)
	})
	
	// 路由重定向：重定向到具体的路由
	r.GET("/redirect2", func(c *gin.Context) {
		c.Request.URL.Path="/TestRedirect"
		r.HandleContext(c)
	})
	// 路由：http://localhost:9090/TestRedirect
	r.GET("/TestRedirect", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"msg":"响应成功",
		})
	})
	r.Run(":9090")
}

```



## 四、返回第三方获取的数据

在我们自己开发的 Server 程序时，Client 请求时需要获取第三方网站上的数据并且将其放回。

#### 核心代码

```go
// 请求第三方数据
reponse, err := http.Get(url)

// 获取响应体
body := response.Body

// 数据返回 Clien
c.DataFromReader(http.StatusOK, contenLength, contentType, body, extraHeaders)

```

Deme

```go
/*
* @Time ： 2022-12-12 23:03
* @Auth ： 张齐林
* @File ：Redirection.go
* @IDE ：GoLand
 */
package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.GET("/GetOtherData", func(c *gin.Context) {
		// url:="https://www.baidu.com"
		url:="https://avatars.githubusercontent.com/u/53826118?v=4"
		response, err := http.Get(url)
		if err != nil || response.StatusCode!= http.StatusOK{
			c.Status(http.StatusServiceUnavailable)  // 应答 Client
			return 
		}
		body := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		// 数据写入响应体
		c.DataFromReader(http.StatusOK,contentLength,contentType,body,nil)
	})
	
	r.Run(":9090")
}

```

## 五、多形式渲染

Server 返回 Client 的数据，需要使用到 `Json` 、`HTML`、`XML`  `YAML` 等多种形式 

#### 核心代码

```go
// 返回 JSON
c.JSON(http.StatusOK,gin.H{"html":"<b>Hello Gin</b>"})

// 返回输出 HTML
c.PureJSON(http.StatusOK,gin.H{"html":"<b>Hello Gin</b>"})

// 返回 YAML 形式(YAML 渲染)
c.YAML(http.StatusOK,gin.H{"message":"hey","status":http.StatusOK})

// 输出 XML 形式
c.XML(http.StatusOK,data)
```

Demo

```go
/*
* @Time ： 2022-12-12 23:03
* @Auth ： 张齐林
* @File ：Redirection.go
* @IDE ：GoLand
 */
package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	// JSON 格式输出
	r.GET("/json", func(c *gin.Context) {
		
		c.JSON(http.StatusOK,gin.H{
			"html":"<b>Hello Gin</b>",
		})
		
	})
	// HTML 格式输出
	r.GET("/someHTML", func(c *gin.Context) {
		c.PureJSON(http.StatusOK,gin.H{
			"html":"<b>Hello Gin</b>",
		})
	})
	
	// XML 格式输出(XML 渲染)
	r.GET("/someXML", func(c *gin.Context) {
		type Message struct {
			Name string
			Msg string
			Age int
		}
		info := Message{
			Name: "张齐林",
			Msg:  "Hello",
			Age:  88,
		}
		c.XML(http.StatusOK,info)
	})
	
	// YAML 格式输出(YAML 渲染)
	r.GET("/someYAML", func(c  *gin.Context) {
		c.YAML(http.StatusOK,gin.H{
			"message":"Gin 框架的多形式渲染",
			"status":http.StatusOK,
		})
	})
	
	r.Run(":9090")
}
```

## 六、文件服务器

Clien 请求的内容是`视频`、`音频`、`图片` 等文件

Gin 框架提供了快速搭建文件服务

#### 核心代码

```go
func(c *gin.Context) {
		path := "File_Path"
		fileName := path+c.Query("name")
		c.File(fileName)
}
```

Demo

```go
/*
* @Time ： 2022-12-15 23:51
* @Auth ： 张齐林
* @File ：File_service.go
* @IDE ：GoLand
 */
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/file", fileServer)
	r.Run(":9090")
}

func fileServer(c *gin.Context) {
	path := "D:/种子/github.com/gin-demo/"
	fileName := path + c.Query("name")
	c.File(fileName)
}
```

## 七、单文件上传

#### 核心代码

```go
file,_ := c.FormFile("file")
// 文件对应的 Key (Post 方法)
c.SaveUploadedFile(file,dst + file.Filename)
// 存储文件
c.String(http.StatusOK,fmt.Sprintf("'%s' 上传完成...",file.Filename))
```

Demo

```go
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
```



## 八、多文件上传

#### 核心代码

```go
form, err := c.MultipartForm()
		if err != nil {
			log.Fatalln(err)
		}
		files := form.File["file_key"]
		for _, file := range files{
			c.SaveUploadedFile(file,dst + file.Filename)
			fmt.Printf("文件 %s 上传成功...\n", file.Filename)
		}
				c.String(http.StatusOK,fmt.Sprintf("%d files uploaded...",len(files)))
```

Demo

```go
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
```

## 九、自定义中间件

Client 请求的路由数据进行预处理（数据加载、请求验证(过滤)等 …）

**中间件里面有错误如果不想继续后续接口的调用不能直接`return`，而是应该调用`c.Abort()`方法。**

源码如下：

```go
// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}
```

核心代码

```go
r ：=gin.New()
r.Use(Middleware)  // 使用中间件
r.Get()
// 自定义路由中间件
func Middleware()gin.HandlerFunc {
	return func(c *gin.Context) {
		
}
```

Demo

```go
/*
* @Time ： 2022-12-17 0:14
* @Auth ： 张齐林
* @File ：Custom_middleware.go
* @IDE ：GoLand
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 如果需要用到没有任何中间件的路由引擎时可以使用以下的 gin.New() 方法
	// r:=gin.New()
	r.Use(Midleware())  // 使用User的话是属于全局中间件的定义使用，如果只需要当单独的方法使用中间件的话(也就是说局部中间件)，直接在方法后面的路由添加即可，如下
	// r.GET("/middleware",Midleware(), func(c *gin.Context) {
	r.GET("/middleware", func(c *gin.Context) {
	fmt.Println("服务端开始执行...")
		
		name := c.Query("name")
		ageStr := c.Query("age")
		age, _ := strconv.Atoi(ageStr)
		log.Println(name, age)
		
		res := struct {
			Name string `json:"name"`
			Age int `json:"age"`
		}{name,age,}
		c.JSON(http.StatusOK,res)
		
	})
	
	r.Run()
}

func Midleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("中间件开始执行...")
		name := c.Query("name")
		ageStr := c.Query("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,"用户数据的数据有错误，age 的值不是整数")
			
			// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
			// Let's say you have an authorization middleware that validates that the current request is authorized.
			// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
			// for this request are not called.
			c.Abort()
			return
		}
		if age <0 || age> 100{
			c.AbortWithStatusJSON(http.StatusBadRequest,"用户数据的数据有错误，age 的值不是整数")
			c.Abort()
			return
		}
		if len(name) <6 || len(name) >12{
			c.AbortWithStatusJSON(http.StatusBadRequest,"用户名只能是6-12位数")
			c.Abort()
			return
		}
		// 执行后续的操作
		c.Next()
		fmt.Println(name, age)
	}
}
```





















