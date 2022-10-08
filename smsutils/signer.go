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
}

/*
{
    "content" : "【线上线下】您的验证码为123456，在10分钟内有效。",
    "mobile" : "13800001111,8613955556666,+8613545556666",
	"extCode" : "123456",
	"sId" : "123456789abcdefg"
}
*/

//短信单发
type SingleSendRequestBody struct {
	//必须。最大长度1005。短信内容，例如： 【线上线下】您的验证码为123456，在10分钟内有效。
	Content string `json:"content"`
	//必须。多个号码用,分隔开，号码数量<=10000
	Mobile string `json:"mobile"`
	//可选。可选。扩展码，必须可解析为数字,最大长度12
	ExtCode string `json:"extCode"`
	//可选。最大长度64。批次号，可用于客户侧按照批次号对短信进行分组
	SId string `json:"sId"`
}

//短信群发，单条短信的结构
type BatchSendItemReuqestBody struct {
	//必须。最大长度1005。短信内容，例如： 【线上线下】您的验证码为123456，在10分钟内有效。
	Content string `json:"content"`
	//必须。多个号码用,分隔开，号码数量<=10000
	Mobile string `json:"mobile"`
	//可选。可选。扩展码，必须可解析为数字,最大长度12
	ExtCode string `json:"extCode"`
	//自定义msgId，若不传，由我们平台生成一个msgId返回，若设置此值，平台将使用此msgId作为此次提交的唯一编号并返回此msgId
	MsgId string `json:"msgid"`
}

type ActiveFetchRequestBody struct {
	//可选 默认500，支持范围[10, 1000]，参数超出范围按照默认算
	MaxSize int `json:"maxSize"`
}

type DailyStatsRequestBody struct {
	//必选,最大长度8，日期格式化：yyyyMMdd 示例： 20200101
	Date string `json:"date"`
}

type SecureSendRequestBody struct {
	Content string `json:"content"`
}

func NewSmsSigner(spId, spKey, smsSendUrl, reportUrl string) *SmsSigner {
	return &SmsSigner{
		SpId:       spId,
		SpKey:      spKey,
		SmsSendUrl: smsSendUrl,
		ReportUrl:  reportUrl,
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
