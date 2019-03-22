package submail

import (
	"fmt"
	"net/url"

	"github.com/json-iterator/go"
)

type SmsXSend struct {
	appId    string
	appKey   string
	signType string

	to      string
	project string
	vars    map[string]string
	tag     string
}

func NewSmsXsend(config map[string]string) *SmsXSend {
	return &SmsXSend{
		appId:    config["appid"],
		appKey:   config["appkey"],
		signType: config["signType"],
		vars:     make(map[string]string),
	}
}

func (this *SmsXSend) SetTo(to string) {
	this.to = to
}

func (this *SmsXSend) SetProject(project string) {
	this.project = project
}

func (this *SmsXSend) AddVar(key string, val string) {
	this.vars[key] = val
}

func (this *SmsXSend) SetTag(tag string) {
	this.tag = tag
}

func (this *SmsXSend) Xsend() (string, error) {
	config := make(map[string]string)

	config["appid"] = this.appId
	config["appkey"] = this.appKey
	config["signType"] = this.signType

	params := url.Values{}

	params.Set("appid", this.appId)
	params.Set("to", this.to)
	params.Set("project", this.project)

	if this.signType != "normal" {
		timestamp, err := GetTimestamp()

		if err != nil {
			return "", err
		}

		params.Set("sign_type", this.signType)
		params.Set("timestamp", fmt.Sprintf("%d", timestamp))
		params.Set("sign_version", "2")
	}

	if this.tag != "" {
		params.Set("tag", this.tag)
	}

	signature := caculSign(params, config)

	params.Set("signature", signature)

	//v2 数字签名 vars 不参与计算
	if len(this.vars) > 0 {
		data, err := jsoniter.Marshal(this.vars)

		if err == nil {
			params.Set("vars", string(data))
		}
	}

	return HttpPost(SUBMAIL_SMS_XSEND_URL, params.Encode())
}
