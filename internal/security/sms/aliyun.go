package sms

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/pkg/errors"
)

type AliyunSMS struct {
	region, accessKey, accessKeySecret string

	signName, templateCode string
}

type NewAliyunSMSOptions struct {
	Region, AccessKey, AccessKeySecret string
	SignName, TemplateCode             string
}

func NewAliyunSMS(options NewAliyunSMSOptions) *AliyunSMS {
	return &AliyunSMS{
		region:          options.Region,
		accessKey:       options.AccessKey,
		accessKeySecret: options.AccessKeySecret,
		signName:        options.SignName,
		templateCode:    options.TemplateCode,
	}
}

func (s *AliyunSMS) SendCode(_ context.Context, phone, code string) error {
	client, err := dysmsapi.NewClientWithAccessKey(s.region, s.accessKey, s.accessKeySecret)
	if err != nil {
		return errors.Wrap(err, "new client")
	}

	req := requests.NewCommonRequest()
	req.Method = http.MethodPost
	req.Scheme = "https"
	req.Domain = "dysmsapi.aliyuncs.com"
	req.Version = "2017-05-25"
	req.ApiName = "SendSms"
	req.QueryParams["RegionId"] = s.region
	req.QueryParams["PhoneNumbers"] = phone
	req.QueryParams["SignName"] = s.signName
	req.QueryParams["TemplateCode"] = s.templateCode
	req.QueryParams["TemplateParam"] = `{"code":"` + code + `"}`
	resp, err := client.ProcessCommonRequest(req)
	if err != nil {
		return errors.Wrap(err, "process common request")
	}

	if err := client.DoAction(req, resp); err != nil {
		return errors.Wrap(err, "do action")
	}

	var result struct {
		Message   string
		Code      string
		BizId     string
		RequestId string
	}
	if err := json.Unmarshal(resp.GetHttpContentBytes(), &result); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	// FYI: https://help.aliyun.com/document_detail/101346.html
	if result.Code != "OK" {
		return errors.Errorf("failed to send sms, code=%v, message=%v, requestId=%v", result.Code, result.Message, result.RequestId)
	}
	return nil
}
