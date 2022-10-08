package smsutils

import "strings"

/**
 * AES 密钥长度规定为16位，这里判断用户名 长度不满足16位的统一前置填充字符'a',大于长度的统一截取前16位
 */

const (
	fillLetter string = "a"
)

func NormalizeKey(key string) string {
	if len(key) >= 16 {
		return key[0:16]
	}
	length := len(key)
	fillLength := 16 - length
	return strings.Repeat(fillLetter, fillLength) + key
}
