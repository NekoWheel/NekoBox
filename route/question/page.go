// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package question

import (
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/mail"
)

func Pager(ctx context.Context) {
	domain := ctx.Param("domain")

	pageUser, err := db.Users.GetByDomain(ctx.Request().Context(), domain)
	if err != nil {
		if errors.Is(err, db.ErrUserNotExists) {
			ctx.Redirect("/")
			return
		} else {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get user by domain")
			ctx.SetError(errors.New("服务器错误！"))
		}
		ctx.Success("question/page")
		return
	}
	ctx.Map(pageUser)

	pageQuestions, err := db.Questions.GetByUserID(ctx.Request().Context(), pageUser.ID, true)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get questions by page id")
		ctx.SetError(errors.New("服务器错误！"))
		ctx.Success("question/page")
		return
	}

	ctx.Data["IsOwnPage"] = ctx.IsLogged && ctx.User.ID == pageUser.ID
	ctx.Data["PageUser"] = pageUser
	ctx.Data["PageQuestions"] = pageQuestions
}

func List(ctx context.Context) {
	ctx.Success("question/list")
}

func New(ctx context.Context, f form.NewQuestion, pageUser *db.User, recaptcha recaptcha.RecaptchaV2) {
	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
		ctx.SetErrorFlash("内部错误，请稍后再试")
		ctx.Redirect("/_/" + pageUser.Domain)
		return
	}
	if !resp.Success {
		ctx.SetErrorFlash("验证码错误")
		ctx.Redirect("/_/" + pageUser.Domain)
		return
	}

	if ctx.HasError() {
		ctx.Success("question/list")
		return
	}

	fromIP := ctx.Request().Header.Get("X-Real-IP")
	question, err := db.Questions.Create(ctx.Request().Context(), db.CreateQuestionOptions{
		FromIP:  fromIP,
		UserID:  pageUser.ID,
		Content: f.Content,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create new question")
		ctx.SetError(errors.New("服务器错误！"))
		ctx.Success("question/list")
		return
	}

	go func() {
		if pageUser.Notify == db.NotifyTypeEmail {
			// Send notification to page user.
			if err := mail.SendNewQuestionMail(pageUser.Email, pageUser.Domain, question.ID, question.Content); err != nil {
				logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to send new question mail to user")
			}
		}
	}()

	ctx.SetSuccessFlash("发送问题成功！")
	ctx.Redirect("/_/" + pageUser.Domain)
}
