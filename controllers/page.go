package controllers

import "github.com/astaxie/beego"

type PageController struct {
	beego.Controller
}

func (this *PageController) Prepare() {
	this.Data["title"] = beego.AppConfig.String("title")
}

// Index is the main page of user's question box.
func (this *PageController) Index() {

}

// NewQuestion is post new question handler.
func (this *PageController) NewQuestion() {

}

// Question is the page of a question.
func (this *PageController) Question() {

}

// AnswerQuestion is the answer question handler.
func (this *PageController) AnswerQuestion() {

}
