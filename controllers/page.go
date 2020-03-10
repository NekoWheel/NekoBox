package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/wuhan005/QuestionBox/models"
	"html/template"
)

type PageController struct {
	beego.Controller
}

func (this *PageController) Prepare() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["icp"] = beego.AppConfig.String("icp")
	this.Data["recaptcha"] = beego.AppConfig.String("recaptcha_site_key")
	this.Data["recaptcha_domain"] = beego.AppConfig.String("recaptcha_domain")
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
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

	// check if the domain is existed.
	domain := this.Ctx.Input.Param(":domain")
	pageContent, err := models.GetPageByDomain(domain)
	if err != nil {
		this.Redirect("/", 302)
		this.Abort("302")
		return
	}
	this.Data["pageContent"] = pageContent
	this.Ctx.Input.SetData("pageContent", pageContent)

	// get the owner of this box
	userContent, _ := models.GetUserByPage(pageContent.ID)
	this.Data["userContent"] = userContent
	this.Ctx.Input.SetData("userContent", userContent)

	// get answer question
	questionContent := models.GetQuestionsByPageID(pageContent.ID, false)
	this.Data["questionContent"] = questionContent
}

// Index is the main page of user's question box.
func (this *PageController) Index() {
	userContent := this.Ctx.Input.GetData("userContent").(*models.User)
	this.Data["title"] = fmt.Sprintf("%s的提问箱 | %s", userContent.Name, beego.AppConfig.String("title"))
}

// NewQuestion is post new question handler.
func (this *PageController) NewQuestion() {
	q := new(models.QuestionForm)
	if err := this.ParseForm(q); err != nil {
		this.Data["error"] = "发送问题失败！"
		this.Data["content"] = q.Content
		return
	}

	valid := validation.Validation{}
	b, err := valid.Valid(q)
	if err != nil {
		this.Data["error"] = "发送问题失败！"
		this.Data["content"] = q.Content
		return
	}
	if !b {
		for _, value := range valid.Errors {
			this.Data["error"] = "问题内容" + value.Message
			this.Data["content"] = q.Content
			return
		}
	}

	// recaptcha
	if !models.CheckRecaptcha(q.Recaptcha, this.Ctx.Input.IP()) {
		this.Data["error"] = "请不要搞事情，感谢。"
		this.Data["content"] = q.Content
		return
	}

	page := this.Ctx.Input.GetData("pageContent").(*models.Page)
	q.PageID = page.ID
	err = models.NewQuestion(q)
	if err != nil {
		this.Data["error"] = err.Error()
		this.Data["content"] = q.Content
		return
	}
	this.Data["success"] = "发送问题成功！"
}
