package util

import (
	"math/rand"
	"time"
)

// 生成随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 获得指定范围随机数字
func GetRandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
