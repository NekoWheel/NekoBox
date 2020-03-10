package controllers

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

func (this *ErrorController) Error404() {
	this.Redirect("/", 302)
	this.TplName = "index.tpl"
}
