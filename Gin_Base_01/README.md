# Gin自学笔记

## 总体学习流程

![image-20221212211824336](README.assets/image-20221212211824336.png)

## 一、Get 请求

HTTP 是超文本传输协议

其定义了客户端和服务端之间文本传输的规范

### Get 方法

Get 方法：主要用于从指定资源中请求数据

使用 Get 方法的请求应该是只是检索数据，并且不应对数据产生其它影响

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





