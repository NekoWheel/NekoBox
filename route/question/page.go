// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package question

import (
	"fmt"

	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wuhan005/govalid"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/dbutil"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/mail"
	"github.com/NekoWheel/NekoBox/internal/security/censor"
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
			ctx.SetInternalError()
		}
		ctx.Success("question/page")
		return
	}
	ctx.Map(pageUser)

	pageQuestions, err := db.Questions.GetByUserID(ctx.Request().Context(), pageUser.ID, db.GetQuestionsByUserIDOptions{
		Cursor:         &dbutil.Cursor{},
		FilterAnswered: true,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get questions by page id")
		ctx.SetInternalError()
		ctx.Success("question/page")
		return
	}

	answeredCount, err := db.Questions.Count(ctx.Request().Context(), pageUser.ID, db.GetQuestionsCountOptions{
		FilterAnswered: true,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to count questions")
		ctx.SetInternalError()
		ctx.Success("question/page")
		return
	}

	ctx.SetTitle(fmt.Sprintf("%s的提问箱 - NekoBox", pageUser.Name))

	ctx.Data["IsOwnPage"] = ctx.IsLogged && ctx.User.ID == pageUser.ID
	ctx.Data["PageUser"] = pageUser
	ctx.Data["PageQuestions"] = pageQuestions
	ctx.Data["CanAsk"] = ctx.IsLogged || pageUser.HarassmentSetting != db.HarassmentSettingTypeRegisterOnly
	ctx.Data["AnsweredCount"] = answeredCount
	if len(pageQuestions) > 0 {
		ctx.Data["PageQuestionCursor"] = pageQuestions[len(pageQuestions)-1].ID
	}
}

func List(ctx context.Context) {
	ctx.Success("question/list")
}

func ListAPI(ctx context.Context) error {
	domain := ctx.Param("domain")
	pageSize := ctx.QueryInt("page_size")
	cursorValue := ctx.Query("cursor")

	pageUser, err := db.Users.GetByDomain(ctx.Request().Context(), domain)
	if err != nil {
		if errors.Is(err, db.ErrUserNotExists) {
			return ctx.JSONError(40400, "用户不存在")
		}
		return ctx.ServerError()
	}

	pageQuestions, err := db.Questions.GetByUserID(ctx.Request().Context(), pageUser.ID, db.GetQuestionsByUserIDOptions{
		Cursor: &dbutil.Cursor{
			Value:    cursorValue,
			PageSize: pageSize,
		},
		FilterAnswered: true,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get questions by page id")
		return ctx.ServerError()
	}

	return ctx.JSON(pageQuestions)
}

func New(ctx context.Context, f form.NewQuestion, pageUser *db.User, recaptcha recaptcha.RecaptchaV3) {
	if !ctx.IsLogged && pageUser.HarassmentSetting == db.HarassmentSettingTypeRegisterOnly {
		ctx.SetErrorFlash("提问箱的主人设置了仅注册用户才能提问，请先登录。")
		ctx.Redirect(fmt.Sprintf("/login?to=%s", ctx.Request().Request.RequestURI))
		return
	}

	var receiveReplyEmail string
	if f.ReceiveReplyViaEmail != "" {
		// Check the email address is valid.
		if errs, ok := govalid.Check(struct {
			Email string `valid:"required;email" label:"邮箱地址"`
		}{
			Email: f.ReceiveReplyEmail,
		}); !ok {
			ctx.SetError(errs[0], f)
			ctx.Success("question/list")
			return
		}

		receiveReplyEmail = f.ReceiveReplyEmail
	}

	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
		ctx.SetInternalErrorFlash()
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

	content := f.Content

	// 🚨 Content security check.
	censorResponse, err := censor.Text(ctx.Request().Context(), content)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to censor text")
	}
	if err == nil && !censorResponse.Pass {
		errorMessage := censorResponse.ErrorMessage()
		ctx.SetError(errors.New(errorMessage), f)
		ctx.Success("question/list")
		return
	}

	// ⚠️ Here is the aliyun CDN origin IP header.
	// A security problem may occur if the CDN is enabled and users can modify the header.
	fromIP := ctx.Request().Header.Get("Ali-CDN-Real-IP")
	if fromIP == "" {
		fromIP = ctx.Request().Header.Get("X-Real-IP")
	}

	// Try to get current logged user.
	var askerUserID uint
	if ctx.IsLogged {
		askerUserID = ctx.User.ID
	}

	question, err := db.Questions.Create(ctx.Request().Context(), db.CreateQuestionOptions{
		FromIP:            fromIP,
		UserID:            pageUser.ID,
		Content:           content,
		ReceiveReplyEmail: receiveReplyEmail,
		AskerUserID:       askerUserID,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create new question")
		ctx.SetInternalError(f)
		ctx.Success("question/list")
		return
	}

	// Update censor result.
	if err := db.Questions.UpdateCensor(ctx.Request().Context(), question.ID, db.UpdateQuestionCensorOptions{
		ContentCensorMetadata: censorResponse.ToJSON(),
	}); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to update question censor result")
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
