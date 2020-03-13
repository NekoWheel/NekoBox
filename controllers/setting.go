package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/wuhan005/QuestionBox/models"
)

type SettingController struct {
	beego.Controller
}

func (this *SettingController) Prepare() {
	isLogin := this.Ctx.Input.GetData("isLogin").(bool)
	if !isLogin {
		this.Redirect("/login", 302)
		this.Abort("302")
		return
	}
}

// Index is the main page of setting.
func (this *SettingController) Index() {
	this.TplName = "setting.tpl"
}

// UpdateProfile is the user's profile update handler.
func (this *SettingController) UpdateProfile() {
	this.TplName = "setting.tpl"
	loginUser := this.Ctx.Input.GetData("user").(*models.User)

	u := new(models.UpdateForm)
	if err := this.ParseForm(u); err != nil {
		this.Data["error"] = "修改个人信息失败！"
		return
	}

	valid := validation.Validation{}
	b, err := valid.Valid(u)
	if err != nil {
		this.Data["error"] = "修改个人信息失败！"
		return
	}
	if !b {
		for _, value := range valid.Errors {
			this.Data["error"] = value.Message
			return
		}
	}
	// check password length
	if len(u.Password) != 0 && len(u.Password) < 8 {
		this.Data["error"] = "密码长度应大于 8 位"
		return
	} else if len(u.Password) > 30 {
		this.Data["error"] = "密码长度应小于 30 位"
		return
	}

	user := &models.User{
		Name: u.Name,
	}
	page := &models.Page{
		Intro: u.Intro,
	}

	// handler picture file
	file, header, err := this.GetFile("avatar")
	if err == nil {
		user.Avatar = models.UploadPicture(header, file)
	}

	file, header, err = this.GetFile("background")
	if err == nil {
		page.Background = models.UploadPicture(header, file)
	}

	// password
	if u.Password != "" {
		user.Password = models.AddSalt(u.Password)
	}

	models.UpdateUser(loginUser.ID, user)
	models.UpdatePage(loginUser.PageID, page)
	this.Data["success"] = "修改个人信息成功！"
}

func (this *SettingController) Logout() {
	this.DestroySession()
	this.Redirect("/", 302)
}
