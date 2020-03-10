package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/parnurzeal/gorequest"
	"github.com/wuhan005/QuestionBox/models"
	"html/template"
)

type SettingController struct {
	beego.Controller
}

func (this *SettingController) Prepare() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["icp"] = beego.AppConfig.String("icp")
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.Data["error"] = ""
	this.Data["success"] = ""

	userInterface := this.GetSession("user")
	if userInterface == nil {
		this.Redirect("/login", 302)
		this.Abort("302")
		return
	}
	user := userInterface.(*models.User)
	this.Data["isLogin"] = true
	this.Data["user"] = user
	this.Ctx.Input.SetData("user", user)
	userPage, _ := models.GetPageByID(user.ID)
	this.Data["page"] = userPage
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
			field := ""
			switch value.Field {
			case "Name":
				field = "昵称"
			case "Password":
				field = "密码"
			case "Intro":
				field = "提问箱介绍"
			}
			this.Data["error"] = field + value.Message
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
		fileByte := make([]byte, header.Size)
		_, _ = file.Read(fileByte)
		req := gorequest.New().Post(beego.AppConfig.String("upload_url")).Type("multipart")
		req.Header.Set("token", beego.AppConfig.String("upload_token"))
		req.SendFile(fileByte, header.Filename, "image")
		resp, body, _ := req.End()

		if resp != nil && resp.StatusCode == 200 {
			avatarJSON := new(models.UploadCallBack)
			err = json.Unmarshal([]byte(body), &avatarJSON)
			if err == nil {
				user.Avatar = avatarJSON.Data.URL
			}
		}
	}

	file, header, err = this.GetFile("background")
	if err == nil {
		fileByte := make([]byte, header.Size)
		_, _ = file.Read(fileByte)
		req := gorequest.New().Post(beego.AppConfig.String("upload_url")).Type("multipart")
		req.Header.Set("token", beego.AppConfig.String("upload_token"))
		req.SendFile(fileByte, header.Filename, "image")
		resp, body, _ := req.End()

		if resp != nil && resp.StatusCode == 200 {
			backgroundJSON := new(models.UploadCallBack)
			err = json.Unmarshal([]byte(body), &backgroundJSON)
			if err == nil {
				page.Background = backgroundJSON.Data.URL
			}
		}
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
