package smsutils

import (
	"net/http"
	"strconv"
	"time"
)

type SmsSigner struct {
	//我方提供的发送账号的唯一标识
	SpId string `json:"spId"`
	//共享spKey密钥
	SpKey string `json:"spKey"`
	//SmsSendUrl 我方提供的短信发送Url地址，例如https://api-hw.onmsg.cn
	SmsSendUrl string `json:"BaseUrl"`
	//平台支持客户主动来获取短信的状态报告，我方提供主动获取短信报告的接口
	ReportUrl string `json:"reportUrl"`
	//模板操作的Url地址
	TemplateUrl string `json:"templateUrl"`
}

/*
{
    "content" : "【线上线下】您的验证码为123456，在10分钟内有效。",
    "mobile" : "13800001111,8613955556666,+8613545556666",
	"extCode" : "123456",
	"sId" : "123456789abcdefg"
}
*/

func NewSmsSigner(spId, spKey, smsSendUrl, reportUrl string, templateUrl string) *SmsSigner {
	return &SmsSigner{
		SpId:        spId,
		SpKey:       spKey,
		SmsSendUrl:  smsSendUrl,
		ReportUrl:   reportUrl,
		TemplateUrl: templateUrl,
	}
}

func (ss *SmsSigner) Sign(r *http.Request, jsonByte []byte) {
	currentTimestring := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	signature := HmacSha256AndBase64(jsonByte, []byte(currentTimestring), []byte(ss.SpKey))
	signatureHeader := "HMAC-SHA256" + " " + currentTimestring + "," + signature
	r.Header.Add("Content-Type", "application/json;charset=utf-8")
	r.Header.Add("Authorization", signatureHeader)
}

// func (ss *SmsSigner) SignWithTimeAgain(jsonByte []byte, timeStamp string) {
// 	signature := HmacSha256AndBase64(jsonByte, []byte(timeStamp), []byte(ss.SpKey))
// 	signatureHeader := "HMAC-SHA256" + " " + timeStamp + "," + signature
// 	fmt.Println("signatureHeader")
// 	fmt.Println(signatureHeader)
// }

// func (ss *SmsSigner) SignAgain(jsonByte []byte) {
// 	//signature := HmacSha256AndBase64(jsonByte,[]byte{},[]byte(ss.SpKey))
// 	h := hmac.New(sha256.New, []byte(ss.SpKey))
// 	h.Write(jsonByte)
// 	buf := h.Sum(nil)
// 	//fmt.Println("sign=" + base64.RawURLEncoding.EncodeToString(buf))
// 	//fmt.Println("sign=" + base64.StdEncoding.EncodeToString(buf))
// 	x := base64.StdEncoding.EncodeToString(buf)
// 	fmt.Println(x)
// }
