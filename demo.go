package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"v4sms/smsutils"
)

//请按实际给的参数填写
const (
	spId       string = "httpSendUser06"
	spKey      string = "httpSendUser06"
	smsSendUrl string = "http://112.25.74.234:38080/"
	reportUrl  string = "https://report-hw.onmsg.cn/"
)

var smsSigner *smsutils.SmsSigner = smsutils.NewSmsSigner(spId, spKey, smsSendUrl, reportUrl)

func main() {
	// normalKey := smsutils.NormalizeKey(spKey)
	// fmt.Println("normalKey is", normalKey)

	fmt.Println("短信发送接口（单内容多号码）")
	demoSingleSend()
	fmt.Println("短信加密发送接口")
	demoSingleSecureSend()
	fmt.Println("短信多发接口")
	demoMultiSend()
	//
	fmt.Println("状态报告主动获取")
	demoStatusFetch()
	fmt.Println("上行主动获取")
	demoUpstreamFetch()
	fmt.Println("预付费账号余额查询")
	demoBalanceFetch()
	fmt.Println("获取发送账号spId的每日短信发送情况统计")
	demoDailyStatsFetch()
	// body := `{"content":"【中秋旅游】综合各地文化和旅游部门、通讯运营商、线上旅行服务商数据，经文化和旅游部数据中心测算，2021年中秋节假期3天，全国累计国内旅游出游8815.93万人次","mobile":"13800001111,13955556666,13545556666","extCode":"123456","sId":"123456789abcdefg"}`
	// smsSigner.SignWithTimeAgain([]byte(body),"1632456539810")
}

func demoSingleSend() {
	requestBody := &smsutils.SingleSendRequestBody{
		"【线上线下】您的验证码为123456，在10分钟内有效。",
		"13800001111,13955556666,13545556666",
		"123456",
		"123456789abcdefg"}

	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("json.Marshal error", err.Error())
		return
	}
	r, _ := http.NewRequest("POST", smsSigner.SingleSendUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func demoSingleSecureSend() {
	requestBody := &smsutils.SingleSendRequestBody{
		"【线上线下】您的验证码为123456，在10分钟内有效。",
		"13800001111,13955556666,13545556666",
		"123456",
		"123456789abcdefg"}

	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("requestBody json.Marshal error", err.Error())
		return
	}

	bAfterEncrypt, err := smsutils.AesECBEncrypt([]byte(jsonByte), []byte(smsutils.NormalizeKey(smsSigner.SpKey)))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	contentString := base64.StdEncoding.EncodeToString(bAfterEncrypt)
	ssrb := &smsutils.SecureSendRequestBody{Content: contentString}
	ssrbByte, err := json.Marshal(ssrb)
	if err != nil {
		fmt.Println("ssrb json.Marshal error", err.Error())
		return
	}
	r, _ := http.NewRequest("POST", smsSigner.SingleSecureSendUrl(), bytes.NewReader(ssrbByte))
	smsSigner.Sign(r, ssrbByte)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func demoMultiSend() {
	smsBody1 := &smsutils.BatchSendItemReuqestBody{
		"【线上线下】线上线下欢迎你参观1",
		"13800001111,8613955556666,+8613545556666",
		"123456",
		"123456787"}

	smsBody2 := &smsutils.BatchSendItemReuqestBody{
		"【线上线下】线上线下欢迎你参观2",
		"13800001111,8613955556666,+8613545556666",
		"123456",
		"123456788"}

	smsBody3 := &smsutils.BatchSendItemReuqestBody{
		"【线上线下】线上线下欢迎你参观3",
		"13800001111,8613955556666,+8613545556666",
		"123456",
		"123456789"}

	requestBody := []*smsutils.BatchSendItemReuqestBody{}
	requestBody = append(requestBody, smsBody1)
	requestBody = append(requestBody, smsBody2)
	requestBody = append(requestBody, smsBody3)
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("json.Marshal error", err.Error())
		return
	}
	//fmt.Println(string(jsonByte))
	r, _ := http.NewRequest("POST", smsSigner.MultiSendUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func demoStatusFetch() {
	requestBody := &smsutils.ActiveFetchRequestBody{500}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("json.Marshal error", err.Error())
		return
	}
	r, _ := http.NewRequest("POST", smsSigner.StatusFetchUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func demoUpstreamFetch() {
	requestBody := &smsutils.ActiveFetchRequestBody{500}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("json.Marshal error", err.Error())
		return
	}
	r, _ := http.NewRequest("POST", smsSigner.UpstreamFetchUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func demoBalanceFetch() {
	r, _ := http.NewRequest("POST", smsSigner.BalanceFetchUrl(), nil)
	smsSigner.Sign(r, nil)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func demoDailyStatsFetch() {
	requestBody := &smsutils.DailyStatsRequestBody{"20200125"}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("json.Marshal error", err.Error())
		return
	}
	r, _ := http.NewRequest("POST", smsSigner.DailyStatsUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
