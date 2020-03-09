package main

import (
	"github.com/astaxie/beego"
	_ "github.com/wuhan005/QuestionBox/routers"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}
