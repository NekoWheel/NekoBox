package routers

import (
	"github.com/astaxie/beego"
	"github.com/wuhan005/QuestionBox/controllers"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
