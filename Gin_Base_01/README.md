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

十、登录中间件

快速实现登录验证中间件方法

Gin 框架提供了快速登录验证中间件，可以完成登录的验证

#### 核心代码

```go
accounts := gin.Accounts{ // gin.Account是map[string]string 类型
		"admin": "adminpw",
	}
// 动态添加用户
accounts["Golang"] = "123654"
accounts["Gin"] = "789abc"
auth := gin.BasicAuth(accounts)
```

Demo

```go
/*
* @Time ： 2022-12-18 23:58
* @Auth ： 张齐林
* @File ：Login_middleware.go
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

	r.Use(AuthMiddleware())
	
	r.GET("/login", func(context *gin.Context) {
		// 获取用户，它是由 BasicAuth 中间件设置的
		user:= context.MustGet(gin.AuthUserKey).(string)
		context.JSON(http.StatusOK,"登录成功..." + "欢迎您: " + user)
	})
	
	err := r.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func AuthMiddleware() gin.HandlerFunc {
	// 初始化用户
	accounts := gin.Accounts{ // gin.Accounts 是 map[string]string
		"admin":"adminpw",
		"system":"systempw",
	}
	
	// 动态添加用户
	accounts["root"]="rootpw"
	
	// 将用户添加到登录中间件中
	auth := gin.BasicAuth(accounts)
	return auth
}
```



## 十、登录中间件

快速实现登录验证中间件方法

Gin 框架提供了快速登录验证中间件，可以完成登录的验证

#### 核心代码

```go
accounts := gin.Accounts{ // gin.Account是map[string]string 类型
		"admin": "adminpw",
	}
// 动态添加用户
accounts["Golang"] = "123654"
accounts["Gin"] = "789abc"
auth := gin.BasicAuth(accounts)
```

Demo

```go
/*
* @Time ： 2022-12-18 23:58
* @Auth ： 张齐林
* @File ：Login_middleware.go
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

	r.Use(AuthMiddleware())
	
	r.GET("/login", func(context *gin.Context) {
		// 获取用户，它是由 BasicAuth 中间件设置的
		user:= context.MustGet(gin.AuthUserKey).(string)
		context.JSON(http.StatusOK,"登录成功..." + "欢迎您: " + user)
	})
	
	err := r.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func AuthMiddleware() gin.HandlerFunc {
	// 初始化用户
	accounts := gin.Accounts{ // gin.Accounts 是 map[string]string
		"admin":"adminpw",
		"system":"systempw",
	}
	
	// 动态添加用户
	accounts["root"]="rootpw"
	
	// 将用户添加到登录中间件中
	auth := gin.BasicAuth(accounts)
	return auth
}
```

## 十一、同步异步

`同步方法`：调用一旦开始，调用者必须等到方法调用返回后，才能继续后续的行为。

`异步方法`：调用更像一个消息队列，一旦开始，方法调用就会立即返回，调用者就可以继续后面的操作，而异步操作通常会在另一个 Go 程（Goroutine）过程，不会障碍调用者的工作。

可以在中间件或处理程序中启动新的Go程（Goroutines）

**特别注意：需要使用上下文的副本**

#### 核心代码

```go
for j:=0;j<6;j++ {
    go func(i int) { // TODO：创建的Go程
        // 创建要在 goroutine 中使用的副本
        cCp := c.Copy()
        time.Sleep(time.Duration(i) * time.Second)
        // 这里使用创建的副本（保证不相互影响）
        fmt.Println("第" + strconv.Itoa(i) + "个Go程: " + cCp.Request.URL.Path)
    }(j)
}
c.JSON(200,"主程序（主Go程）" + c.Request.URL.Path)
```

Demo

```go
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
```



## 十二、多服务器程序运行

程序运行在多个服务器上

多服务器程序同时运行，最先想到的方法是就是：写多个 Gin 框架对应的程序

其实，Gin 框架提供了快捷的方法	

#### 核心代码

```go
func main() {
    // 服务器1：http://127.0.0.1:8080/MulServer01
	server01 := &http.Server{
		Addr: ":8080",
		Handler: router01(),
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
	// 服务器2：http://127.0.0.1:8081/MulServer02
	server02 := &http.Server{
		Addr: ":8081",
		Handler: router02(),
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
}
```

Demo

```go
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
```

## 十三、路由组

一个项目中可以有多个路由，同时可以通过路由组来继续管理和有效的分类，使路由对应的代码易于阅读

#### 核心代码

```go
router := gin.Default()
	// 路由分组1
	v1 := router.Group("/v1")  // 路由分组1(一级路径)
	{
		r := v1.Group("/user")     // 路由分组(二级路径)
		r.GET("/login",login)      // 响应请求: v1/user/login
		// 路由分组(三级路径)
		r2 := r.Group("/showInfo")  // /v1/user/showInfo
		r2.GET("/abstract",abstract)   // 响应请求: v1/user/showInfo/abstract
		r2.GET("detail",detail)   // 响应请求：v1/user/showInfo/detail
		
	}
	// 路由分组2
	v2 := router.Group("/v2")
	{
		v2.GET("/other",other)
	}
```

Demo

```go
/*
* @Time ： 2022-12-21 23:43
* @Auth ： 张齐林
* @File ：Routing_group.go
* @IDE ：GoLand
 */
package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type ResGroup struct {
	Data string
	Path string
}

func main() {
	
	router := gin.Default()
	// 路由分组1
	v1 := router.Group("/v1")  // 路由分组1(一级路径)
	{
		r := v1.Group("/user")     // 路由分组(二级路径)
		r.GET("/login",login)      // 响应请求: v1/user/login
		// 路由分组(三级路径)
		r2 := r.Group("/showInfo")  // /v1/user/showInfo
		r2.GET("/abstract",abstract)   // 响应请求: v1/user/showInfo/abstract
		r2.GET("detail",detail)   // 响应请求：v1/user/showInfo/detail
		
	}
	// 路由分组2
	v2 := router.Group("/v2")
	{
		v2.GET("/other",other)
	}
	
	err := router.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func other(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "other",
		Path: c.Request.URL.Path,
	})
}

func detail(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "detail",
		Path: c.Request.URL.Path,
	})
}

func abstract(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "abstract",
		Path: c.Request.URL.Path,
	})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK,ResGroup{
		Data: "login",
		Path: c.Request.URL.Path,
	})
} 
```

## 十四、Gin 框架Bind

将 Client 提交的 JSON 数据与 Server 对应的对象（实体/结构体）进行关联。

Gin 框架提供了 Bind ，可以根据请求Body数据，将数据赋值到指定的结构体变量中（类似于序列化和反序列化）

Gin 的 Bind 方法，主要是将结构体与请求参数进行绑定，请求参数 JSON 对应的 Key 就是结构体对应的字段。

#### 核心代码

```go
type Login struct {
	User string`form:"user"binding:"required"`
	PassWord string`form:"password"binding:"required,min=6,max=12"`
}

c.Bing(&login)
```

Demo

```go
/*
* @Time ： 2022-12-22 22:43
* @Auth ： 张齐林
* @File ：Gin_Bind.go
* @IDE ：GoLand
 */
package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type Login struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remark string `json:"remark"`
}

func main() {
	
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		
		var login Login
		err := c.Bind(&login)
		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"msg":"绑定失败,参数错误",
				"data":err.Error(),
			})
			return
		}
		if login.UserName == "user" && login.Password == "123456"{
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":"登录成功...",
			"data":"OK",
		})
			return
		}
		c.JSON(http.StatusBadRequest,gin.H{
			"msg":"登录失败...",
			"data":"error",
		})
	})
	
	err := r.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
```

## 十五、调用 Restful 接口

Server 程序需要调用其它接口（Restful接口）

#### 核心代码

```go
// data             Post 请求提交的数据
// url              请求地址
// contentTyper     请求体格式，如：application/json
// contemt          接口返回的内容
func Post(url string,data any,contentType string)string  {
	client := &http.Client{Timeout: time.Second * 5}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url,contentType,bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	// result, _ := ioutil.ReadAll(resp.Body)
	result, _ := io.ReadAll(resp.Body)
	return string(result)
}
```

Demo

```go
/*
* @Time ： 2022-12-23 22:31
* @Auth ： 张齐林
* @File ：Call_restfu_interface.go
* @IDE ：GoLand
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
)

// UserAPI 调用第三方接口的请求数据结构体
type UserAPI struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Other string `json:"other"`
}

// TempData 调用第三方接口的返回结果
type TempData struct {
	Msg string `json:"msg"`
	Data string `json:"data"`
}

// ClientRequest 客户端提交的数据
type ClientRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Other string `json:"other"`
}

// ClientResponse 返回客户端的数据
type ClientResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
}

// data             Post 请求提交的数据
// url              请求地址
// contentTyper     请求体格式，如：application/json
// content          接口返回的内容
func getRestfulAPI(url string,data any, contentType  string)([]byte,error)  {
	client := &http.Client{Timeout: time.Second * 5}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType,bytes.NewBuffer(jsonStr))  // client.Get() 的Url中一定要携带传递的参数
	if err != nil {
		fmt.Println("调用API接口出现了错误...")
		return nil, err
	}
	defer resp.Body.Close()
	// result, _ := ioutil.ReadAll(resp.Body)
	result, err := io.ReadAll(resp.Body)
	return result,err
}

func testAPI()  {
	url := "http://127.0.0.1:8080/login"
	user := UserAPI{"user","123456",""}
	data, err := getRestfulAPI(url,user,"application/json")
	fmt.Println(data,err)
	var temp TempData
	json.Unmarshal(data,&temp)
	fmt.Println(temp.Msg, temp.Data)
}

func main() {
	// testAPI()
	r := gin.Default()
	r.POST("/getOtherAPI",getOtherAPI)
	r.Run(":9090")
}

func getOtherAPI(context *gin.Context) {
	var requestData ClientRequest
	var response ClientResponse
	err := context.Bind(&requestData)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Msg = "请求的参数错误..."
		response.Data= err
		context.JSON(http.StatusBadRequest,response)
		return
	}
	// 请求第三方数据
	url := "http://127.0.0.1:8080/login"
	user := UserAPI{requestData.UserName,requestData.Password,requestData.Other}
	data, err := getRestfulAPI(url,user,"application/json")
	fmt.Println(data,err)
	var temp TempData
	json.Unmarshal(data,&temp)
	fmt.Println(temp.Msg, temp.Data)
	response.Code = http.StatusOK
	response.Msg = "请求成功..."
	response.Data = temp
	context.JSON(http.StatusOK,response)
}
```

## 十六、参数校验_1

由于GoLang没有跟 Java 一样的注解快速进行参数的校验，但是可以通过struct tag （结构体标签）进行序列化

如以下常用的

```go
type User struct {
    ID string `josn:"id"`
    Name string `json:"name"`
    Age string `json:"age"`
}
```

如果需要更复杂的校验，这时需要用一些专业的库来完成

go-playground / validator 作为一款优秀（使用简单、快捷）的 Go 语言校验库，基于标记作为结构体和单个字段实现值校验。

<u>Gin 框架在 binding 中已嵌套了</u>

#### 核心代码

```go
```

Demo

```go
```



## 十七、参数校验_2



#### 核心代码

```go

```

Demo

```go
```



## 十八、

## 十九、

## 二十、

## 二十一、

















































































