package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

const (
	// 随机字符串源
	randStrSource = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)


// 获取随机字符串
func RanString(l int)string{
	length := len(randStrSource)
	r := make([]byte, l)
	for i := 0; i < l; i++{
		r[i] = randStrSource[rand.Intn(length)]
	}
	return string(r)
}

// 计算字符串的md5
func Md5(s string)string{
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
