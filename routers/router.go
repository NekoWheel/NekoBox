package routers

import (
	"html/template"

	"github.com/NekoWheel/NekoBox/controllers"
	"github.com/NekoWheel/NekoBox/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

const (
	COMMITSHA string = "aef6d648f7ca9923550c504ac2407e9d979e0734"
)

func init() {
	beego.InsertFilter("*", beego.BeforeExec, func(c *context.Context) {
		c.Input.Data()["title"] = beego.AppConfig.String("title")
		c.Input.Data()["icp"] = beego.AppConfig.String("icp")
		c.Input.Data()["commitSha"] = COMMITSHA
		c.Input.Data()["recaptcha"] = beego.AppConfig.String("recaptcha_site_key")
		c.Input.Data()["recaptcha_domain"] = beego.AppConfig.String("recaptcha_domain")
		c.Input.Data()["xsrfdata"] = template.HTML(`<input type="hidden" name="_xsrf" value="` +
			c.XSRFToken(beego.BConfig.WebConfig.XSRFKey, int64(beego.BConfig.WebConfig.XSRFExpire)) +
			`" />`)
		c.Input.Data()["success"] = ""
		c.Input.Data()["error"] = ""

		// get login status
		user := c.Input.Session("user")
		if user != nil {
			c.Input.Data()["isLogin"] = true
			c.Input.Data()["user"] = user.(*models.User)
			c.Input.SetData("user", user.(*models.User))
			c.Input.SetData("isLogin", true)

			userPage, _ := models.GetPageByID(user.(*models.User).ID)
			c.Input.Data()["page"] = userPage
		} else {
			c.Input.Data()["isLogin"] = false
			c.Input.SetData("isLogin", false)
		}
	})

	beego.Router("/", &controllers.MainController{})

	beego.Router("/register", &controllers.UserController{}, "get:RegisterGet;post:RegisterPost")
	beego.Router("/login", &controllers.UserController{}, "get:LoginGet;post:LoginPost")

	beego.Router("/_/:domain", &controllers.PageController{}, "get:Index;post:NewQuestion")
	beego.Router("/_/:domain/:id:int", &controllers.QuestionController{}, "get:Question;post:AnswerQuestion")

	beego.Router("/question", &controllers.QuestionController{}, "get:QuestionList")
	beego.Router("/delete/:domain/:id:int", &controllers.QuestionController{}, "post:QuestionDelete")
	beego.Router("/setting", &controllers.SettingController{}, "get:Index;post:UpdateProfile")
	beego.Router("/logout", &controllers.SettingController{}, "get:Logout")

	beego.ErrorController(&controllers.ErrorController{})
}
