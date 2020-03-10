package routers

import (
	"github.com/astaxie/beego"
	"github.com/wuhan005/QuestionBox/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/register", &controllers.UserController{}, "get:RegisterGet;post:RegisterPost")
	beego.Router("/login", &controllers.UserController{}, "get:LoginGet;post:LoginPost")

	beego.Router("/_/:domain", &controllers.PageController{}, "get:Index;post:NewQuestion")
	beego.Router("/_/:domain/:id:int", &controllers.QuestionController{}, "get:Question;post:AnswerQuestion")

	beego.Router("/question", &controllers.QuestionController{}, "get:QuestionList")
	beego.Router("/setting", &controllers.SettingController{}, "get:Index;post:UpdateProfile")
	beego.Router("/logout", &controllers.SettingController{}, "get:Logout")

}
