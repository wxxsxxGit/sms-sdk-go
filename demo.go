package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"v4sms/pkg/strutils"
	"v4sms/smsutils"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var globalTemplateCode int64
var smsSigner *smsutils.SmsSigner

// var smsSigner *smsutils.SmsSigner = smsutils.NewSmsSigner(spId, spKey, smsSendUrl, reportUrl, templateUrl)

func configPrompt() {
	fmt.Println("配置文件默认为/etc/sms.yaml\n" +
		"需要5个配置项spId,spKey,smsSendUrl,reportUrl,templateUrl联系管理员获取")
}

func init() {
	viper.SetConfigName("sms")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc") // call multiple times to add many search paths
	viper.AddConfigPath(".")    // path to look for the config file in
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("../../config/")
	viper.AddConfigPath("../../../config/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(color.RedString("config file not found:" + err.Error()))
			configPrompt()
			os.Exit(111)
		} else {
			fmt.Println(color.RedString("config file  found,can not use:" + err.Error()))
			configPrompt()
			os.Exit(112)
		}
	}
	spId := viper.GetString("spId")
	spKey := viper.GetString("spKey")
	smsSendUrl := viper.GetString("smsSendUrl")
	reportUrl := viper.GetString("reportUrl")
	templateUrl := viper.GetString("templateUrl")
	if len(spId) == 0 ||
		len(spKey) == 0 ||
		len(smsSendUrl) == 0 ||
		len(reportUrl) == 0 ||
		len(templateUrl) == 0 {
		configPrompt()
		os.Exit(113)
	}
	smsSigner = smsutils.NewSmsSigner(spId, spKey, smsSendUrl, reportUrl, templateUrl)

}

func main() {

	//单条内容发送
	fmt.Println("短信发送接口（单内容多号码）")
	demoSingleSend()
	sperator(1)

	//单条内容加密发送
	fmt.Println("短信加密发送接口")
	demoSingleSecureSend()
	sperator(1)

	//多内容批量发送
	fmt.Println("短信多发接口")
	demoMultiSend()
	sperator(1)

	//主动获取状态报告
	fmt.Println("状态报告主动获取")
	demoStatusFetch()
	sperator(1)

	//主动获取上行
	fmt.Println("上行主动获取")
	demoUpstreamFetch()
	sperator(1)

	//查询余额
	fmt.Println("预付费账号余额查询")
	demoBalanceFetch()
	sperator(1)

	//查询每日发送统计
	fmt.Println("获取发送账号spId的每日短信发送情况统计")
	demoDailyStatsFetch()
	sperator(1)

	//模板报备
	fmt.Println("模板报备")
	templateId, err := demoTemplateAdd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sperator(1)
	//把服务端生成的templateCode设置为公共的值
	globalTemplateCode = templateId
	fmt.Println("模板报备的code为", globalTemplateCode)

	fmt.Println("等待模板审核...")
	//判断是否审核成功
	var tempValue uint8
	for {
		m, err := demoTemplateStatus(globalTemplateCode)
		if err != nil {
			fmt.Println("demoTemplateStatus()", err.Error())
			time.Sleep(5 * time.Second)
		}
		value, ok := m[globalTemplateCode]
		if !ok {
			fmt.Println("出现错误", globalTemplateCode, "不存在")
			return
		}
		if value == 0 {
			fmt.Println(globalTemplateCode, "联系管理员审核", time.Now().String())
			time.Sleep(10 * time.Second)
			continue
		} else {
			tempValue = value
			break
		}
	}

	//审核成功提交模板短信
	if tempValue == 1 {
		fmt.Println(globalTemplateCode, "审核通过")
		sperator(1)
	} else if tempValue == 2 {
		//审核失败修改模板，只有在模板审核失败时才可以修改模板
		fmt.Println("模板审核失败,模板修改后提交")
		err = demoTemplateModify(globalTemplateCode)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//模板被禁用
	} else {
		fmt.Println(globalTemplateCode, "审核状态为", tempValue, "退出")
	}

	//模板发送单条短信
	fmt.Println("模板发送单条短信")
	err = demoTemplateSendSms(globalTemplateCode)
	if err != nil {
		fmt.Println("demoTemplateSendSms", err.Error())
		return
	}
	sperator(1)

	//模板批量发送短信
	fmt.Println("模板发送批量短信")
	err = demoTemplateSendBatchSms(globalTemplateCode)
	if err != nil {
		fmt.Println("demoTemplateSendBatchSms", err.Error())
		return
	}
	sperator(1)

	//删除模板
	fmt.Println("10秒后将删除模板", globalTemplateCode)
	time.Sleep(10 * time.Second)
	err = demoTemplateDelete(globalTemplateCode)
	if err != nil {
		fmt.Println("demoTemplateDelete", err.Error())
		return
	}
	fmt.Println("删除模板", globalTemplateCode, "成功")
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

//模板报备http请求demo
func demoTemplateAdd() (int64, error) {
	requestBody := &smsutils.TemplateAddRequestBody{

		TemplateName:    "线上线下addTemplate SDK DEMO " + strutils.RandString(10),
		TemplateType:    2,
		TemplateContent: "线上线下addTemplate SDK DEMO template content ${code} template " + strutils.RandString(20),
		Remark:          "线上线下addTemplate SDK DEMO template " + time.Now().Format("2006-01-02 15:04:05"),
	}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}
	r, _ := http.NewRequest("POST", smsSigner.TemplateAddUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	fmt.Println(string(body))
	tarb := &smsutils.TemplateAddRespBody{}
	err = json.Unmarshal(body, &tarb)
	if err != nil {
		return 0, err
	}
	return tarb.TemplateCode, nil
}

//模板修改 http请求demo
func demoTemplateModify(templateCode int64) error {
	requestBody := &smsutils.TemplateModifyRequestBody{
		TemplateCode:    templateCode,
		TemplateType:    2,
		TemplateName:    "线上线下addTemplate SDK DEMO " + strutils.RandString(10),
		TemplateContent: "线上线下addTemplate SDK DEMO template content ${code} modify template " + strutils.RandString(20),
		Remark:          "线上线下addTemplate SDK DEMO template " + time.Now().Format("2006-01-02 15:04:05"),
	}

	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	r, _ := http.NewRequest("POST", smsSigner.TemplateModifyUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	tmrb := &smsutils.TemplateModifyRespBody{}
	err = json.Unmarshal(body, &tmrb)
	if err != nil {
		return err
	}
	if tmrb.Status != 0 {
		return errors.New("status为" + strconv.FormatInt(int64(tmrb.Status), 10) + ",msg为" + tmrb.Msg)
	}
	return nil
}

//模板查询 http请求demo
func demoTemplateStatus(templateCodes ...int64) (map[int64]uint8, error) {
	tSlice := []string{}
	for _, v := range templateCodes {
		tSlice = append(tSlice, strconv.FormatInt(v, 10))
	}
	requestBody := &smsutils.TemplateStatusRequestBody{
		TemplateCodes: strings.Join(tSlice, ","),
	}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	r, _ := http.NewRequest("POST", smsSigner.TemplateStatusUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http状态码为" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	tsrb := &smsutils.TemplateStatusRespBody{}
	err = json.Unmarshal(body, &tsrb)
	if err != nil {
		return nil, err
	}
	if tsrb.Status != 0 {
		return nil, errors.New("status为" + strconv.FormatInt(int64(tsrb.Status), 10) + ",msg为" + tsrb.Msg)
	}
	m := make(map[int64]uint8)
	for _, v := range tsrb.TemplateList {
		m[v.TemplateCode] = v.AuditStatus
	}
	return m, nil
}

//模板删除 请求demo
func demoTemplateDelete(templateCode int64) error {
	requestBody := &smsutils.TemplateDeleteRequestBody{
		TemplateCode: templateCode,
	}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	r, _ := http.NewRequest("POST", smsSigner.TemplateDeleteUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("http状态码为" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	tdrb := &smsutils.TemplateDeleteRespBody{}
	err = json.Unmarshal(body, &tdrb)
	if err != nil {
		return err
	}
	if tdrb.Status != 0 {
		return errors.New("status为" + strconv.FormatInt(int64(tdrb.Status), 10) + ",msg为" + tdrb.Msg)
	}

	return nil
}

//模板单条发送 请求demo
func demoTemplateSendSms(templateCode int64) error {
	params := make(map[string]string)
	params["code"] = "普通的一个"
	paramsByte, err := json.Marshal(params)
	if err != nil {
		return err
	}
	requestBody := &smsutils.TemplateSendSmsRequestItem{
		SignName:     "模板测试",
		TemplateCode: templateCode,
		Params:       string(paramsByte),
		Mobile:       "18799991367,12899190876,13914117531",
	}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	r, _ := http.NewRequest("POST", smsSigner.TemplateSendSmsUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("http状态码为" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(body))
	tdrb := &smsutils.TemplateSendSmsRespBody{}
	err = json.Unmarshal(body, &tdrb)
	if err != nil {
		return err
	}
	if tdrb.Status != 0 {
		return errors.New("status为" + strconv.FormatInt(int64(tdrb.Status), 10) + ",msg为" + tdrb.Msg)
	}
	fmt.Println("成功号码为", tdrb.SuccessList)
	fmt.Println("失败号码为", tdrb.FailList)
	fmt.Println("短信分片为", tdrb.SplitCount)
	fmt.Println("msgid为", tdrb.Msgid)
	return nil
}

//模板批量发送 请求demo
func demoTemplateSendBatchSms(templateCode int64) error {
	params := make(map[string]string)
	params["code"] = "普通的一个"
	paramsByte, err := json.Marshal(params)
	if err != nil {
		return err
	}
	item1 := &smsutils.TemplateSendSmsRequestItem{
		SignName:     "模板测试1",
		TemplateCode: templateCode,
		Params:       string(paramsByte),
		Mobile:       "18505101387",
	}
	item2 := &smsutils.TemplateSendSmsRequestItem{
		SignName:     "模板测试2",
		TemplateCode: templateCode,
		Params:       string(paramsByte),
		Mobile:       "12899190872",
	}
	item3 := &smsutils.TemplateSendSmsRequestItem{
		SignName:     "模板测试3",
		TemplateCode: templateCode,
		Params:       string(paramsByte),
		Mobile:       "18799991362",
	}
	item4 := &smsutils.TemplateSendSmsRequestItem{
		SignName:     "模板测试4",
		TemplateCode: templateCode,
		Params:       string(paramsByte),
		Mobile:       "13914117532",
	}
	item5 := &smsutils.TemplateSendSmsRequestItem{
		SignName:     "模板测试5",
		TemplateCode: templateCode,
		Params:       string(paramsByte),
		Mobile:       "1895606996",
	}

	requestBody := &smsutils.TemplateSendBatchSmsRequestBody{
		item1,
		item2,
		item3,
		item4,
		item5,
	}
	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	r, _ := http.NewRequest("POST", smsSigner.TemplateSendBatchSmsUrl(), bytes.NewReader(jsonByte))
	smsSigner.Sign(r, jsonByte)

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("http状态码为" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	tdrb := &smsutils.TemplateSendBatchSmsRespBody{}
	err = json.Unmarshal(body, &tdrb)
	if err != nil {
		return err
	}
	if tdrb.Status != 0 {
		return errors.New("status为" + strconv.FormatInt(int64(tdrb.Status), 10) + ",msg为" + tdrb.Msg)
	}
	for _, v := range tdrb.Result {
		fmt.Println(v)
	}
	return nil
}

func sperator(sec int) {
	fmt.Println(strings.Repeat("*", 30) + "\n")
	time.Sleep(time.Duration(sec) * time.Second)
}
