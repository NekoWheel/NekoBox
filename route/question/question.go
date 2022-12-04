// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package question

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/mail"
	"github.com/NekoWheel/NekoBox/internal/security/censor"
)

func Questioner(ctx context.Context, pageUser *db.User) {
	questionID := uint(ctx.ParamInt("questionID"))
	question, err := db.Questions.GetByID(ctx.Request().Context(), questionID)
	if err != nil {
		if !errors.Is(err, db.ErrQuestionNotExist) {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get question by ID")
		}

		ctx.Redirect("/")
		return
	}
	ctx.Data["Question"] = question

	// Check the question is belongs to the correct page user.
	// If the question has not been answered, we should check the question is belongs to the correct page user.
	if question.UserID != pageUser.ID || (question.Answer == "" && (!ctx.IsLogged || ctx.User.ID != question.UserID)) {
		ctx.Redirect("/")
		return
	}

	// The page's owner or the question's token can have the permission to delete the question.
	// Inject the permission into the context.
	token := ctx.Query("t")
	canDelete := (ctx.IsLogged && ctx.User.ID == pageUser.ID) || (token == question.Token && question.Token != "")
	ctx.Map(canDelete)
	ctx.Data["CanDelete"] = canDelete

	ctx.Map(question)
}

func Item(ctx context.Context) {
	ctx.Success("question/item")
}

func PublishAnswer(ctx context.Context, pageUser *db.User, question *db.Question, f form.PublishAnswerQuestion) {
	if ctx.HasError() {
		ctx.Success("question/item")
		return
	}

	if ctx.User.ID != pageUser.ID {
		ctx.Redirect("/")
		return
	}

	answer := f.Answer

	// 🚨 Content security check.
	censorResponse, err := censor.Text(ctx.Request().Context(), answer)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to censor text")
	}
	if err == nil && !censorResponse.Pass {
		errorMessage := censorResponse.ErrorMessage()
		ctx.SetError(errors.New(errorMessage), f)
		ctx.Success("question/item")
		return
	}

	if err := db.Questions.AnswerByID(ctx.Request().Context(), question.ID, f.Answer); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to answer question")
		ctx.SetError(errors.New("服务器错误！"), f)
		ctx.Success("question/item")
		return
	}

	// Update censor result.
	if err := db.Questions.UpdateCensor(ctx.Request().Context(), question.ID, db.UpdateQuestionCensorOptions{
		AnswerCensorMetadata: censorResponse.ToJSON(),
	}); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to update answer censor result")
	}

	go func() {
		if question.ReceiveReplyEmail != "" && question.Answer == "" { // We only send the email when the question has not been answered.
			// Send notification to questioner.
			if err := mail.SendNewAnswerMail(pageUser.Email, pageUser.Domain, question.ID, question.Content, f.Answer); err != nil {
				logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to send receive reply mail to questioner")
			}
		}
	}()

	ctx.SetSuccessFlash("回答发布成功！")
	ctx.Redirect(fmt.Sprintf("/_/%s/%d", pageUser.Domain, question.ID))
}

func Delete(ctx context.Context, pageUser *db.User, question *db.Question, canDelete bool) {
	if !canDelete {
		ctx.Redirect("/_/" + pageUser.Domain)
		return
	}

	if err := db.Questions.DeleteByID(ctx.Request().Context(), question.ID); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to delete question")
		ctx.SetError(errors.New("服务器错误！"))
		ctx.Success("question/item")
		return
	}

	ctx.Redirect("/_/" + pageUser.Domain)
}
