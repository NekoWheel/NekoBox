// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"encoding/gob"
	"time"

	"github.com/flamego/cache"
	"github.com/flamego/csrf"
	"github.com/flamego/flamego"
	"github.com/flamego/recaptcha"
	"github.com/flamego/session"
	"github.com/flamego/session/mysql"
	"github.com/flamego/template"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/form"
	templatepkg "github.com/NekoWheel/NekoBox/internal/template"
	"github.com/NekoWheel/NekoBox/route"
	"github.com/NekoWheel/NekoBox/route/auth"
	"github.com/NekoWheel/NekoBox/route/question"
	"github.com/NekoWheel/NekoBox/route/user"
	"github.com/NekoWheel/NekoBox/templates"
)

func New() *flamego.Flame {
	f := flamego.Classic()
	if conf.App.Production {
		flamego.SetEnv(flamego.EnvTypeProd)
	}

	templateFS, err := template.EmbedFS(templates.FS, ".", []string{".html"})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to embed templates file system")
	}

	// We prefer to save session into database,
	// if no database configuration, the session will be saved into memory instead.
	gob.Register(time.Time{})
	gob.Register(context.Flash{})
	var sessionStorage interface{}
	initer := session.MemoryIniter()
	if conf.Database.DSN != "" {
		initer = mysql.Initer()
		sessionStorage = mysql.Config{
			DSN:      conf.Database.DSN,
			Lifetime: 7 * 24 * time.Hour,
		}
	}
	sessioner := session.Sessioner(session.Options{
		Initer: initer,
		Config: sessionStorage,
	})

	reqUserSignOut := context.Toggle(&context.ToggleOptions{UserSignOutRequired: true})
	reqUserSignIn := context.Toggle(&context.ToggleOptions{UserSignInRequired: true})

	f.Group("", func() {
		f.Get("/", route.Home)
		f.Get("/sponsor", route.Sponsor)
		f.Get("/change-logs", route.ChangeLogs)
		f.Get("/robots.txt", func(c context.Context) {
			_, _ = c.ResponseWriter().Write([]byte("User-agent: *\nDisallow: /_/"))
		})

		f.Group("", func() {
			f.Combo("/register").Get(auth.Register).Post(form.Bind(form.Register{}), auth.RegisterAction)
			f.Combo("/login").Get(auth.Login).Post(form.Bind(form.Login{}), auth.LoginAction)
			f.Combo("/forgot-password").Get(auth.ForgotPassword).Post(form.Bind(form.ForgotPassword{}), auth.ForgotPasswordAction)
			f.Combo("/recover-password").Get(auth.RecoverPassword).Post(form.Bind(form.RecoverPassword{}), auth.RecoverPasswordAction)
		}, reqUserSignOut)

		f.Group("/_/{domain}", func() {
			f.Combo("").Get(question.List).Post(form.Bind(form.NewQuestion{}), question.New)
			f.Group("/{questionID}", func() {
				f.Get("", question.Item)
				//f.Post("/update", question.Update)
				f.Post("/delete", question.Delete)
				f.Post("/answer", reqUserSignIn, form.Bind(form.PublishAnswerQuestion{}), question.PublishAnswer)
			}, question.Questioner)
		}, question.Pager)

		f.Group("/user", func() {
			f.Get("/questions", user.QuestionList)

			f.Group("/profile", func() {
				f.Get("", user.Profile)
				f.Post("/update", form.Bind(form.UpdateProfile{}), user.UpdateProfile)
				f.Post("/export", user.ExportProfile)
				f.Combo("/deactivate").Get(user.DeactivateProfile).Post(user.DeactivateProfileAction)
			})

			f.Get("/logout", auth.Logout)
		}, reqUserSignIn)
	},
		cache.Cacher(),
		recaptcha.V2(
			recaptcha.Options{
				Secret:    conf.Recaptcha.ServerKey,
				VerifyURL: recaptcha.VerifyURLGlobal,
			},
		),
		sessioner,
		csrf.Csrfer(csrf.Options{
			Secret: conf.Server.XSRFKey,
			Header: "X-CSRF-Token",
		}),
		template.Templater(template.Options{
			FileSystem: templateFS,
			FuncMaps:   templatepkg.FuncMap(),
		}),
		context.Contexter(),
	)
	f.NotFound(func(ctx flamego.Context) {
		ctx.Redirect("/")
	})

	return f
}
