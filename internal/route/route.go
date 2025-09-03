// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"
	"strings"
	"time"

	"github.com/flamego/cache"
	cacheRedis "github.com/flamego/cache/redis"
	"github.com/flamego/flamego"
	"github.com/flamego/recaptcha"
	"github.com/flamego/session"
	"github.com/flamego/session/mysql"
	sessionRedis "github.com/flamego/session/redis"
	"gorm.io/gorm"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/form"
)

func New(db *gorm.DB) *flamego.Flame {
	f := flamego.Classic()
	if conf.App.Production {
		flamego.SetEnv(flamego.EnvTypeProd)
	}

	// We prefer to save session into database,
	// if no database configuration, the session will be saved into memory instead.
	var sessionStorage interface{}
	initer := session.MemoryIniter()
	if conf.Database.DSN != "" {
		initer = mysql.Initer()
		sessionStorage = mysql.Config{
			DSN:      conf.Database.DSN,
			Lifetime: 7 * 24 * time.Hour,
		}
	}
	if conf.Redis.Addr != "" {
		initer = sessionRedis.Initer()
		sessionStorage = sessionRedis.Config{
			Options: &cacheRedis.Options{
				Addr:     conf.Redis.Addr,
				Password: conf.Redis.Password,
				DB:       1,
			},
			Lifetime: 30 * 24 * time.Hour,
		}
	}
	sessioner := session.Sessioner(session.Options{
		Initer: initer,
		Config: sessionStorage,
		ReadIDFunc: func(r *http.Request) string {
			authorizationHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			groups := strings.Fields(authorizationHeader)
			if len(groups) != 2 {
				return ""
			}
			return groups[1]
		},
		WriteIDFunc: func(_ http.ResponseWriter, _ *http.Request, _ string, _ bool) {},
	})

	f.Use(
		sessioner,
		cache.Cacher(cache.Options{
			Initer: cacheRedis.Initer(),
			Config: cacheRedis.Config{
				Options: &cacheRedis.Options{
					Addr:     conf.Redis.Addr,
					Password: conf.Redis.Password,
					DB:       0,
				},
			},
		}),
		recaptcha.V3(
			recaptcha.Options{
				Secret: conf.Recaptcha.ServerKey,
				VerifyURL: func() recaptcha.VerifyURL {
					if conf.Recaptcha.TurnstileStyle {
						// FYI: https://developers.cloudflare.com/turnstile/migration/migrating-from-recaptcha/
						return "https://challenges.cloudflare.com/turnstile/v0/siteverify"
					}
					return recaptcha.VerifyURLGlobal
				}(),
			},
		),
		context.Contexter(db),
	)

	reqUserSignOut := context.Toggle(&context.ToggleOptions{UserSignOutRequired: true})
	reqUserSignIn := context.Toggle(&context.ToggleOptions{UserSignInRequired: true})

	f.Group("/api", func() {
		authHandler := NewAuthHandler()
		f.Group("/auth", func() {
			f.Post("/sign-up", form.Bind(form.SignUp{}), authHandler.SignUp)
			f.Post("/sign-in", form.Bind(form.SignIn{}), authHandler.SignIn)
			f.Post("/forgot-password", form.Bind(form.ForgotPassword{}), authHandler.ForgotPassword)
			f.Combo("/recover-password").Get(authHandler.GetRecoverPasswordCode).Post(form.Bind(form.RecoverPassword{}), authHandler.RecoverPassword)
		}, reqUserSignOut)

		userHandler := NewUserHandler()
		f.Group("/users/{domain}", func() {
			f.Get("/profile", userHandler.Profile)
			f.Group("/questions", func() {
				f.Combo("").
					Get(userHandler.ListQuestions).
					Post(form.BindMultipart(form.PostQuestion{}), userHandler.PostQuestion)
				f.Get("/{questionID}", userHandler.GetQuestion)
			})
		}, userHandler.Domainer)

		mineHandler := NewMineHandler()
		f.Group("/mine", func() {
			f.Group("/questions", func() {
				f.Get("", mineHandler.ListQuestions)
				f.Group("/{questionID}", func() {
					f.Put("/answer", form.BindMultipart(form.AnswerQuestion{}), mineHandler.AnswerQuestion)
					f.Delete("", mineHandler.DeleteQuestion)
					f.Put("/visible", form.Bind(form.QuestionVisible{}), mineHandler.SetQuestionVisible)
				}, mineHandler.Questioner)
			})

			f.Group("/settings", func() {
				f.Combo("/profile").Get(mineHandler.Profile).Put(form.Bind(form.UpdateProfile{}), mineHandler.UpdateProfileSettings)
				f.Combo("/box").Get(mineHandler.BoxSettings).Put(form.BindMultipart(form.UpdateBoxSettings{}), mineHandler.UpdateBoxSettings)
				f.Combo("/harassment").Get(mineHandler.HarassmentSettings).Put(form.Bind(form.UpdateHarassmentSettings{}), mineHandler.UpdateHarassmentSettings)
				f.Post("/export-data", mineHandler.ExportData)
				f.Post("/deactivate", mineHandler.Deactivate)
			})
		}, reqUserSignIn)
	})

	//f.Group("", func() {
	//	f.Get("/", route.Home)
	//	f.Get("/pixel", reqUserSignIn, pixel.Index)
	//	f.Get("/sponsor", route.Sponsor)
	//	f.Get("/change-logs", route.ChangeLogs)
	//	f.Get("/robots.txt", func(c context.Context) {
	//		_, _ = c.ResponseWriter().Write([]byte("User-agent: *\nDisallow: /_/"))
	//	})
	//	f.Get("/favicon.ico", func(c context.Context) {
	//		fs, _ := static.FS.Open("favicon.ico")
	//		defer func() { _ = fs.Close() }()
	//		c.ResponseWriter().Header().Set("Content-Type", "image/x-icon")
	//		_, _ = io.Copy(c.ResponseWriter(), fs)
	//	})
	//
	//	f.Group("", func() {
	//		f.Combo("/register").Get(auth.Register).Post(form.Bind(form.Register{}), auth.RegisterAction)
	//		f.Combo("/login").Get(auth.Login).Post(form.Bind(form.Login{}), auth.LoginAction)
	//		f.Combo("/forgot-password").Get(auth.ForgotPassword).Post(form.Bind(form.ForgotPassword{}), auth.ForgotPasswordAction)
	//		f.Combo("/recover-password").Get(auth.RecoverPassword).Post(form.Bind(form.RecoverPassword{}), auth.RecoverPasswordAction)
	//	}, reqUserSignOut)
	//
	//	f.Group("/_/{domain}", func() {
	//		f.Combo("").Get(question.List).Post(form.Bind(form.NewQuestion{}), question.New)
	//		f.Group("/{questionID}", func() {
	//			f.Get("", question.Item)
	//			f.Post("/delete", question.Delete)
	//			f.Post("/set-private", question.SetPrivate)
	//			f.Post("/set-public", question.SetPublic)
	//			f.Post("/answer", reqUserSignIn, form.Bind(form.PublishAnswerQuestion{}), question.PublishAnswer)
	//		}, question.Questioner)
	//	}, question.Pager)
	//
	//	f.Group("/user", func() {
	//		f.Get("/questions", user.QuestionList)
	//
	//		f.Group("/profile", func() {
	//			f.Get("", user.Profile)
	//			f.Post("/update", form.Bind(form.UpdateProfile{}), user.UpdateProfile)
	//			f.Post("/export", user.ExportProfile)
	//			f.Combo("/deactivate").Get(user.DeactivateProfile).Post(user.DeactivateProfileAction)
	//		})
	//		f.Post("/harassment/update", form.Bind(form.UpdateHarassment{}), user.UpdateHarassment)
	//
	//		f.Get("/logout", auth.Logout)
	//	}, reqUserSignIn)
	//
	//	f.Group("/api/v1", func() {
	//		f.Group("/user", func() {
	//			f.Get("", reqUserSignIn, user.ProfileAPI)
	//
	//			f.Group("/{domain}", func() {
	//				f.Group("/questions", func() {
	//					f.Get("", question.ListAPI)
	//				})
	//			})
	//		})
	//
	//		f.Any("/pixel/{**}", reqUserSignIn, pixel.Proxy)
	//	}, context.APIEndpoint)
	//},
	//	cache.Cacher(cache.Options{
	//		Initer: cacheRedis.Initer(),
	//		Config: cacheRedis.Config{
	//			Options: &cacheRedis.Options{
	//				Addr:     conf.Redis.Addr,
	//				Password: conf.Redis.Password,
	//				DB:       0,
	//			},
	//		},
	//	}),
	//	recaptcha.V3(
	//		recaptcha.Options{
	//			Secret: conf.Recaptcha.ServerKey,
	//			VerifyURL: func() recaptcha.VerifyURL {
	//				if conf.Recaptcha.TurnstileStyle {
	//					// FYI: https://developers.cloudflare.com/turnstile/migration/migrating-from-recaptcha/
	//					return "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	//				}
	//				return recaptcha.VerifyURLGlobal
	//			}(),
	//		},
	//	),
	//	sessioner,
	//	csrf.Csrfer(csrf.Options{
	//		Secret: conf.Server.XSRFKey,
	//		Header: "X-CSRF-Token",
	//	}),
	//	template.Templater(template.Options{
	//		FileSystem: templateFS,
	//		FuncMaps:   templatepkg.FuncMap(),
	//	}),
	//	context.Contexter(),
	//)
	f.NotFound(func(ctx flamego.Context) {
		ctx.Redirect("/")
	})

	return f
}
