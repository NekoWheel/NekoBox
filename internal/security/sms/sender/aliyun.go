// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sender

import (
	"context"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/pkg/errors"

	"github.com/NekoWheel/NekoBox/internal/security/sms/storage"
)

var _ Sender = (*AliyunSMS)(nil)

type AliyunSMS struct {
	signName                           string
	region, accessKey, accessKeySecret string
}

type NewAliyunSMSOptions struct {
	SignName                           string
	Region, AccessKey, AccessKeySecret string
}

func NewAliyunSMS(options NewAliyunSMSOptions) *AliyunSMS {
	return &AliyunSMS{
		signName:        options.SignName,
		region:          options.Region,
		accessKey:       options.AccessKey,
		accessKeySecret: options.AccessKeySecret,
	}
}

func (s *AliyunSMS) Send(_ context.Context, typ storage.SMSType, phone, code string) error {
	client, err := dysmsapi.NewClientWithAccessKey(s.region, s.accessKey, s.accessKeySecret)
	if err != nil {
		return errors.Wrap(err, "new client")
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = s.signName

	var templateCode, templateParam string
	switch typ {
	case storage.SMSTypeVerifyPhone:
		templateCode = "SMS_261120088" // FIXME: I am not sure if this is stable.
		templateParam = `{"code":"` + code + `"}`
	default:
		return errors.Errorf("unknown sms type: %q", typ)
	}
	request.TemplateCode = templateCode
	request.TemplateParam = templateParam

	response, err := client.SendSms(request)
	if err != nil {
		return errors.Wrap(err, "send sms")
	}

	// https://help.aliyun.com/document_detail/101346.html
	if response.Code != "OK" {
		return errors.Errorf("send sms failed, code: %v, message: %v", response.Code, response.Message)
	}
	return nil
}
