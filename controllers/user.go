package controllers

import (
	"github.com/NekoWheel/NekoBox/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Prepare() {
	isLogin, _ := this.Ctx.Input.GetData("isLogin").(bool)
	if isLogin {
		user := this.Data["user"].(*models.User)
		domain, _ := models.GetPageByID(user.PageID)
		this.Redirect("/_/"+domain.Domain, 302)
		this.Abort("302")
		return
	}
}

// RegisterGet: user register page
func (this *UserController) RegisterGet() {
	this.TplName = "register.tpl"
}

// Post: user register handler
func (this *UserController) RegisterPost() {
	this.TplName = "register.tpl"

	r := new(models.UserRegisterForm)
	if err := this.ParseForm(r); err != nil {
		this.Data["error"] = "注册失败！"
		this.Data["name"] = r.Name
		this.Data["email"] = r.Email
		this.Data["domain"] = r.Domain
		return
	}

	valid := validation.Validation{}
	b, err := valid.Valid(r)
	if err != nil {
		this.Data["error"] = "注册失败！"
		this.Data["name"] = r.Name
		this.Data["email"] = r.Email
		this.Data["domain"] = r.Domain
		return
	}
	if !b {
		for _, value := range valid.Errors {
			this.Data["error"] = value.Message
			this.Data["name"] = r.Name
			this.Data["email"] = r.Email
			this.Data["domain"] = r.Domain
			return
		}
	}

	if !models.CheckRecaptcha(r.Recaptcha, this.Ctx.Input.IP()) {
		this.Data["error"] = "请不要搞事情，感谢。"
		this.Data["name"] = r.Name
		this.Data["email"] = r.Email
		this.Data["domain"] = r.Domain
		return
	}

	err = models.Register(r)
	if err != nil {
		this.Data["error"] = err.Error()
		this.Data["name"] = r.Name
		this.Data["email"] = r.Email
		this.Data["domain"] = r.Domain
		return
	}

	this.Redirect("/login", 302)
}

// LoginGet: user login page
func (this *UserController) LoginGet() {
	this.TplName = "login.tpl"
}

// LoginPost: user login handler
func (this *UserController) LoginPost() {
	this.TplName = "login.tpl"
	r := new(models.UserLoginForm)
	if err := this.ParseForm(r); err != nil {
		this.Data["error"] = "登录失败！"
		this.Data["email"] = r.Email
		return
	}

	valid := validation.Validation{}
	b, err := valid.Valid(r)
	if err != nil {
		this.Data["error"] = "登录失败！"
		this.Data["email"] = r.Email
		return
	}
	if !b {
		for _, value := range valid.Errors {
			this.Data["error"] = value.Message
			this.Data["email"] = r.Email
			return
		}
	}

	// recaptcha
	if !models.CheckRecaptcha(r.Recaptcha, this.Ctx.Input.IP()) {
		this.Data["error"] = "请不要搞事情，感谢。"
		this.Data["email"] = r.Email
		return
	}

	user, err := models.Login(r)
	if err != nil {
		this.Data["error"] = "用户名或密码错误！"
		this.Data["email"] = r.Email
		return
	}

	page, err := models.GetPageByID(user.PageID)
	if err != nil {
		this.Data["error"] = "用户名或密码错误！"
		this.Data["email"] = r.Email
		return
	}

	this.SetSession("user", user)

	this.Redirect("/_/"+page.Domain, 302)
}
