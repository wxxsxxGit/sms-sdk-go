package smsutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

//这里golang跟java有区别，不能直接拼接字符串然后转化成[]byte。必须用下面的方法
//参考https://stackoverflow.com/questions/51259453/hmacsha256-signature-generated-is-different-in-java-than-in-go
func HmacSha256AndBase64(b1 []byte, b2 []byte, spKeyByte []byte) string {
	h := hmac.New(sha256.New, spKeyByte)
	h.Write(b1)
	h.Write(b2)
	buf := h.Sum(nil)
	//fmt.Println("sign=" + base64.RawURLEncoding.EncodeToString(buf))
	//fmt.Println("sign=" + base64.StdEncoding.EncodeToString(buf))
	return base64.StdEncoding.EncodeToString(buf)
}
