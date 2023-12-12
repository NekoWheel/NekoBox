// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"mime/multipart"
)

type NewQuestion struct {
	Content              string `form:"content" valid:"required;maxlen:1000" label:"问题内容"`
	ReceiveReplyViaEmail string
	ReceiveReplyEmail    string                  `label:"接收回复的电子邮箱"`
	Recaptcha            string                  `form:"g-recaptcha-response" valid:"required" label:"Recaptcha" msg:"无感验证码加载错误，请尝试刷新页面重试。"`
	Images               []*multipart.FileHeader `form:"images" label:"图片"`
}

type PublishAnswerQuestion struct {
	Answer string                  `form:"answer" valid:"required;maxlen:1000" label:"回答内容"`
	Images []*multipart.FileHeader `form:"images" label:"图片"`
}

type UpdateAnswerQuestion struct {
	Answer string                  `form:"answer" valid:"required;maxlen:1000" label:"回答内容"`
	Images []*multipart.FileHeader `form:"images" label:"图片"`
}
