package smsutils

import (
	"strings"
)

//短信发送接口（单内容多号码）url
func (ss *SmsSigner) SingleSendUrl() string {
	return strings.TrimRight(ss.SmsSendUrl, "/") + "/sms/send/" + ss.SpId
}

//短信加密发送接口（影响性能，非必要不推荐）url
func (ss *SmsSigner) SingleSecureSendUrl() string {
	return strings.TrimRight(ss.SmsSendUrl, "/") + "/sms/secureSend/" + ss.SpId
}

//多内容发送url
func (ss *SmsSigner) MultiSendUrl() string {
	//fmt.Println(strings.TrimRight(ss.SmsSendUrl,"/")+"/sms/sendBatch/"+ss.SpId)
	return strings.TrimRight(ss.SmsSendUrl, "/") + "/sms/sendBatch/" + ss.SpId
}

//状态报告主动获取url
func (ss *SmsSigner) StatusFetchUrl() string {
	return strings.TrimRight(ss.ReportUrl, "/") + "/sms/getReport/" + ss.SpId
}

//上行主动获取url
func (ss *SmsSigner) UpstreamFetchUrl() string {
	return strings.TrimRight(ss.ReportUrl, "/") + "/sms/getUpstream/" + ss.SpId
}

//预付费账号余额查询url
func (ss *SmsSigner) BalanceFetchUrl() string {
	return strings.TrimRight(ss.ReportUrl, "/") + "/sms/getBalance/" + ss.SpId
}

//获取发送账号日统计url
func (ss *SmsSigner) DailyStatsUrl() string {
	return strings.TrimRight(ss.ReportUrl, "/") + "/sms/getDailyStats/" + ss.SpId
}

//
