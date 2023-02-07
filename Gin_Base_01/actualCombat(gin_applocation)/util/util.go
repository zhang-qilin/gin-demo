/*
* @Time ： 2023-02-06 19:55
* @Auth ： 张齐林
* @File ：util.go
* @IDE ：GoLand
 */
package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMMQWERTYUIO")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
