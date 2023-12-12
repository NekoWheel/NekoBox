// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

type UpdateProfile struct {
	Name        string `valid:"required;maxlen:20" label:"昵称"`
	OldPassword string `label:"旧密码"`
	NewPassword string `valid:"maxlen:30" label:"新密码"`
	Phone       string `label:"手机号码"`
	VerifyCode  string `label:"手机验证码"`
	Intro       string `valid:"required;maxlen:100" label:"介绍"`
	NotifyEmail string `label:"开启邮箱通知"`
}

type UpdateHarassment struct {
	RegisterOnly string `label:"仅允许注册用户"`
}
