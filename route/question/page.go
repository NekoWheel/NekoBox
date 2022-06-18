// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package question

import (
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
)

func Pager(ctx context.Context) {
	domain := ctx.Param("domain")

	pageUser, err := db.Users.GetByDomain(ctx.Request().Context(), domain)
	if err != nil {
		if errors.Is(err, db.ErrUserNotExists) {
			ctx.Redirect("/")
			return
		} else {
			log.Error("Failed to get user by domain: %v", err)
			ctx.SetError(errors.New("服务器错误！"))
		}
		ctx.Success("question/page")
		return
	}
	ctx.Map(pageUser)

	pageQuestions, err := db.Questions.GetByUserID(ctx.Request().Context(), pageUser.ID, true)
	if err != nil {
		log.Error("Failed to get questions by page id: %v", err)
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
		log.Error("Failed to check recaptcha: %v", err)
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

	if err := db.Questions.Create(ctx.Request().Context(), db.CreateQuestionOptions{
		UserID:  pageUser.ID,
		Content: f.Content,
	}); err != nil {
		log.Error("Failed to create new question: %v", err)
		ctx.SetError(errors.New("服务器错误！"))
		ctx.Success("question/list")
		return
	}

	ctx.SetSuccessFlash("发送问题成功！")
	ctx.Redirect("/_/" + pageUser.Domain)
}
