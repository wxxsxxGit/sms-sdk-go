package smsutils

type BaseResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

//模板报备请求body
type TemplateAddRespBody struct {
	TemplateCode int64 `json:"templateCode"`
	BaseResp
}

//修改模板body
type TemplateModifyRespBody struct {
	TemplateCode int64 `json:"templateCode"`
	BaseResp
}

//删除模板body
type TemplateDeleteRespBody struct {
	BaseResp
}

//查询模板状态body
type TemplateStatusRespBody struct {
	BaseResp
	TemplateList []*TemplateStatusItem
}

type TemplateStatusItem struct {
	BaseResp
	TemplateCode    int64  `json:"templateCode"`
	TemplateType    int    `json:"templateType"`
	TemplateContent string `json:"templateContent"`
	TemplateName    string `json:"templateName"`
	AuditStatus     uint8  `json:"auditStatus"`
	AuditReason     string `json:"auditReason"`
	CreateTime      string `json:"createTime"`
}

//模板单条发送body
type TemplateSendSmsRespBody struct {
	BaseResp
	FailList    string `json:"failList"`
	SuccessList string `json:"successList"`
	SplitCount  int    `json:"splitCount"`
	Msgid       string `json:"msgid"`
}

//模板批量发送Body
type TemplateSendBatchSmsItem struct {
	BaseResp
	Mobile     string `json:"mobile"`
	Msgid      string `json:"msgid"`
	SplitCount int    `json:"splitCount"`
}

type TemplateSendBatchSmsRespBody struct {
	BaseResp
	Result []*TemplateSendBatchSmsItem `json:"result"`
}
