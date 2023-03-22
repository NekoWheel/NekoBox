// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

type Register struct {
	Email          string `valid:"required;email;maxlen:100" label:"电子邮箱"`
	Domain         string `valid:"required;alphadash;minlen:3;maxlen:20" label:"个性域名"`
	Phone          string `label:"手机号码"`
	VerifyCode     string `label:"手机验证码"`
	Name           string `valid:"required;maxlen:20" label:"昵称"`
	Password       string `valid:"required;minlen:8;maxlen:30" label:"密码"`
	RepeatPassword string `valid:"required;equal:Password" label:"重复密码"`
	Recaptcha      string `form:"g-recaptcha-response" valid:"required" label:"Recaptcha" msg:"无感验证码加载错误，请尝试刷新页面重试。"`
}

type SendSMS struct {
	Phone     string `valid:"required;phone" label:"手机号码"`
	Recaptcha string `form:"g-recaptcha-response" valid:"required" label:"Recaptcha" msg:"无感验证码加载错误，请尝试刷新页面重试。"`
}

type Login struct {
	Email     string `valid:"required;email;maxlen:100" label:"电子邮箱"`
	Password  string `valid:"required" label:"密码"`
	Recaptcha string `form:"g-recaptcha-response" valid:"required" label:"Recaptcha"`
}

type ForgotPassword struct {
	Email     string `valid:"required;email;maxlen:100" label:"电子邮箱"`
	Recaptcha string `form:"g-recaptcha-response" valid:"required" label:"Recaptcha"`
}

type RecoverPassword struct {
	NewPassword    string `valid:"required;minlen:8;maxlen:30" label:"新密码"`
	RepeatPassword string `valid:"required;equal:NewPassword" label:"重复密码"`
}
