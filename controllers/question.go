package controllers

import (
	"github.com/astaxie/beego"
	"github.com/wuhan005/QuestionBox/models"
	"strconv"
)

type QuestionController struct {
	beego.Controller
}

func (this *QuestionController) Prepare() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["success"] = ""
	this.Data["error"] = ""

	// get login status
	user := this.GetSession("user")
	if user != nil {
		this.Data["isLogin"] = true
		this.Ctx.Input.SetData("isLogin", true)
		this.Data["user"] = user.(*models.User)
		this.Ctx.Input.SetData("user", user.(*models.User))

		userPage, _ := models.GetPageByID(user.(*models.User).ID)
		this.Data["page"] = userPage
	} else {
		this.Data["isLogin"] = false
		this.Data["user"] = nil
		this.Ctx.Input.SetData("isLogin", false)
	}
}

// Question is the page of a question.
func (this *QuestionController) Question() {
	domain := this.Ctx.Input.Param(":domain")
	id := this.Ctx.Input.Param(":id")
	questionID, err := strconv.Atoi(id)
	if err != nil {
		this.Redirect("/", 302)
		return
	}

	question, err := models.GetQuestionByDomainID(domain, uint(questionID))
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	user, _ := models.GetUserByPage(question.PageID)
	page, _ := models.GetPageByDomain(domain)
	questions := models.GetQuestionsByPageID(question.PageID, false)
	this.Data["userContent"] = user
	this.Data["pageContent"] = page
	this.Data["questionsContent"] = questions
	this.Data["questionContent"] = question
	this.TplName = "question.tpl"
}

// QuestionList show the owner's all questions.
func (this *QuestionController) QuestionList() {
	isLogin := this.Ctx.Input.GetData("isLogin").(bool)
	if !isLogin {
		this.Redirect("/login", 302)
		return
	}
	user := this.Ctx.Input.GetData("user").(*models.User)
	questions := models.GetQuestionsByPageID(user.PageID, true)
	this.Data["questionContent"] = questions
	this.TplName = "questionlist.tpl"
}

// AnswerQuestion is the answer question handler.
func (this *QuestionController) AnswerQuestion() {

}
