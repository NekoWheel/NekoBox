package form

import "mime/multipart"

type AnswerQuestion struct {
	Answer string                  `form:"answer" valid:"required;maxlen:1000" label:"回答内容"`
	Images []*multipart.FileHeader `form:"images[]" label:"图片"`
}

type QuestionVisible struct {
	Visible bool `json:"visible"`
}

type UpdateProfile struct {
	Name        string `json:"name" valid:"required;maxlen:20" label:"昵称"`
	OldPassword string `json:"oldPassword" label:"旧密码"`
	NewPassword string `json:"newPassword" valid:"maxlen:30" label:"新密码"`
}
