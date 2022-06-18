// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package question

import (
	"fmt"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
)

func Questioner(ctx context.Context, pageUser *db.User) {
	questionID := uint(ctx.ParamInt("questionID"))
	question, err := db.Questions.GetByID(ctx.Request().Context(), questionID)
	if err != nil {
		if !errors.Is(err, db.ErrQuestionNotExist) {
			log.Error("Failed to get question by ID: %v", err)
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
	canDelete := (ctx.IsLogged && ctx.User.ID == pageUser.ID) || (token == question.Token)
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

	if err := db.Questions.AnswerByID(ctx.Request().Context(), question.ID, f.Answer); err != nil {
		log.Error("Failed to answer question: %v", err)
		ctx.SetError(errors.New("服务器错误！"))
		return
	}

	ctx.SetSuccessFlash("回答发布成功！")
	ctx.Redirect(fmt.Sprintf("/_/%s/%d", pageUser.Domain, question.ID))
}

func Delete(ctx context.Context, pageUser *db.User, question *db.Question, canDelete bool) {
	if !canDelete {
		ctx.Redirect("/_/" + pageUser.Domain)
		return
	}

	if err := db.Questions.DeleteByID(ctx.Request().Context(), question.ID); err != nil {
		log.Error("Failed to delete question: %v", err)
		ctx.SetError(errors.New("服务器错误！"))
		ctx.Success("question/item")
		return
	}

	ctx.Redirect("/_/" + pageUser.Domain)
}
