// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"mime/multipart"
)

type PostQuestion struct {
	Content           string                  `form:"content" valid:"required;maxlen:1000" label:"提问内容"`
	ReceiveReplyEmail string                  `form:"receiveReplyEmail"`
	Images            []*multipart.FileHeader `form:"images[]" label:"图片"`
	IsPrivate         bool                    `form:"isPrivate"`
	Recaptcha         string                  `form:"recaptcha" valid:"required" label:"无感验证码" msg:"无感验证码加载错误，请尝试刷新页面重试。"`
}
