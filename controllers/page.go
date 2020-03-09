package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/wuhan005/QuestionBox/models"
)

type PageController struct {
	beego.Controller
}

func (this *PageController) Prepare() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["success"] = ""
	this.Data["error"] = ""

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

// Index is the main page of user's question box.
func (this *PageController) Index() {
	// check if the domain is existed.
	domain := this.Ctx.Input.Param("domain")
	page, err := models.GetPageByDomain(domain)
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	this.Data["pageContent"] = page

	// get the owner of this box
	user, _ := models.GetUserByPage(page.ID)
	this.Data["userContent"] = user

	this.TplName = "page.tpl"
}

// NewQuestion is post new question handler.
func (this *PageController) NewQuestion() {
	this.TplName = "page.tpl"

	// check if the domain is existed.
	domain := this.Ctx.Input.Param("domain")
	page, err := models.GetPageByDomain(domain)
	if err != nil {
		this.Redirect("/", 302)
		return
	}

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

	q.PageID = page.ID
	err = models.NewQuestion(q)
	if err != nil {
		this.Data["error"] = err.Error()
		this.Data["content"] = q.Content
		return
	}
	this.Data["success"] = "发送问题成功！"
}

// Question is the page of a question.
func (this *PageController) Question() {

}

// AnswerQuestion is the answer question handler.
func (this *PageController) AnswerQuestion() {

}
