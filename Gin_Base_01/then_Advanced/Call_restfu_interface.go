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