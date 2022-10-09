package smsutils

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

//主动获取状态
type ActiveFetchRequestBody struct {
	//可选 默认500，支持范围[10, 1000]，参数超出范围按照默认算
	MaxSize int `json:"maxSize"`
}

//获取日常报告
type DailyStatsRequestBody struct {
	//必选,最大长度8，日期格式化：yyyyMMdd 示例： 20200101
	Date string `json:"date"`
}

//加密短信
type SecureSendRequestBody struct {
	Content string `json:"content"`
}

//模板报备body
type TemplateAddBody struct {
	TemplateName    string `json:"templateName"`
	TemplateType    int    `json:"templateType"`
	TemplateContent string `json:"templateContent"`
	Remark          string `json:"remark"`
}

//修改模板body
type TemplateModifyBody struct {
	TemplateCode    int64  `json:"templateCode"`
	TemplateName    string `json:"templateName"`
	TemplateType    int    `json:"templateType"`
	TemplateContent string `json:"templateContent"`
	Remark          string `json:"remark"`
}

//删除模板body
type TemplateDeleteBody struct {
	TemplateCode int64 `json:"templateCode"`
}

//查询模板状态body
type TemplateStatusBody struct {
	TemplateCodes string `json:"templateCodes"`
}

//模板单条发送body
type TemplateSendSmsBody struct {
	SignName     string            `json:"signName"`
	TemplateCode int64             `json:"templateCode"`
	Params       map[string]string `json:"params"`
	Mobile       string            `json:"mobile"`
	Msgid        string            `json:"msgid"`
	ExtCode      string            `json:"extCode"`
	SId          string            `json:"sId"`
}

//模板批量发送Body
type TemplateSendBatchSmsBody []*TemplateSendSmsBody
