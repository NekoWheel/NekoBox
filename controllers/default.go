package controllers

import (
	"github.com/astaxie/beego"
	"github.com/wuhan005/QuestionBox/models"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Prepare() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["icp"] = beego.AppConfig.String("icp")
	this.Data["success"] = ""
	this.Data["error"] = ""
	this.TplName = "page.tpl"

	// get login status
	user := this.GetSession("user")
	if user != nil {
		this.Data["isLogin"] = true
		this.Data["user"] = user.(*models.User)

		userPage, _ := models.GetPageByID(user.(*models.User).ID)
		this.Data["page"] = userPage
	} else {
		this.Data["isLogin"] = false
	}
}

func (this *MainController) Get() {

	this.TplName = "index.tpl"
}
