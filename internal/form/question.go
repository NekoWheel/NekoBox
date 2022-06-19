// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

type NewQuestion struct {
	Content   string `form:"content" valid:"required;maxlen:1000" label:"问题内容"`
	Recaptcha string `form:"g-recaptcha-response" valid:"required" label:"Recaptcha"`
}

type PublishAnswerQuestion struct {
	Answer string `form:"answer" valid:"required;maxlen:1000" label:"回答内容"`
}

type UpdateAnswerQuestion struct {
	Answer string `form:"answer" valid:"required;maxlen:1000" label:"回答内容"`
}
