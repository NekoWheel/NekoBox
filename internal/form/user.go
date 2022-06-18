// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

type UpdateProfile struct {
	Name        string `valid:"required;maxlen:20" label:"回答内容"`
	OldPassword string `label:"旧密码"`
	NewPassword string `valid:"maxlen:30" label:"新密码"`
	Intro       string `valid:"required;avatar;maxlen:100" label:"介绍"`
}
