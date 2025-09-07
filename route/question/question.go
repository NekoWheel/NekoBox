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
	token := ctx.Query("t")
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

	askUploadImages, err := db.UploadImgaes.GetByTypeQuestionID(ctx.Request().Context(), db.UploadImageQuestionTypeAsk, questionID)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get ask upload images")
	}
	ctx.Data["AskUploadImages"] = askUploadImages

	answerUploadImages, err := db.UploadImgaes.GetByTypeQuestionID(ctx.Request().Context(), db.UploadImageQuestionTypeAnswer, questionID)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get answer upload images")
	}
	ctx.Data["AnswerUploadImages"] = answerUploadImages

	// Check the question is belongs to the correct page user.
	// If the question has not been answered, we should check the question is belongs to the correct page user.
	// The questioner can use the token to view the question.
	if question.UserID != pageUser.ID ||
		((question.Answer == "" || question.IsPrivate) &&
			(!ctx.IsLogged || ctx.User.ID != question.UserID) &&
			(question.Token != "" && question.Token != token)) {
		ctx.Redirect("/")
		return
	}

	// The page's owner or the question's token can have the permission to delete the question.
	// Inject the permission into the context.
	canDelete := ctx.IsLogged && ctx.User.ID == pageUser.ID
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

	if len(f.Images) > 0 {
		image := f.Images[0]

		if err := uploadImage(ctx, uploadImageOptions{
			Type:               db.UploadImageQuestionTypeAnswer,
			Image:              image,
			QuestionID:         question.ID,
			UploaderUserID:     ctx.User.ID,
			IsDeletingPrevious: true,
		}); err != nil {
			if errors.Is(err, ErrUploadImageSizeTooLarge) {
				ctx.SetErrorFlash("图片文件大小不能大于 5Mb")
				ctx.Success("question/item")
				return
			} else {
				logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to upload image")
			}
		}
	}

	if err := db.Questions.AnswerByID(ctx.Request().Context(), question.ID, f.Answer); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to answer question")
		ctx.SetInternalError(f)
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
			if err := mail.SendNewAnswerMail(question.ReceiveReplyEmail, pageUser.Domain, question.ID, question.Content, f.Answer); err != nil {
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
		ctx.SetInternalError()
		ctx.Success("question/item")
		return
	}

	ctx.Redirect("/_/" + pageUser.Domain)
}

func SetPrivate(ctx context.Context, pageUser *db.User, question *db.Question, canDelete bool) {
	if !canDelete {
		ctx.Redirect("/_/" + pageUser.Domain)
		return
	}

	if err := db.Questions.SetPrivate(ctx.Request().Context(), question.ID); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to set question private")
		ctx.SetInternalError()
		ctx.Success("question/item")
		return
	}

	ctx.Redirect(fmt.Sprintf("/_/%s/%d", pageUser.Domain, question.ID))
}

func SetPublic(ctx context.Context, pageUser *db.User, question *db.Question, canDelete bool) {
	if !canDelete {
		ctx.Redirect("/_/" + pageUser.Domain)
		return
	}

	if err := db.Questions.SetPublic(ctx.Request().Context(), question.ID); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to set question public")
		ctx.SetInternalError()
		ctx.Success("question/item")
		return
	}

	ctx.Redirect(fmt.Sprintf("/_/%s/%d", pageUser.Domain, question.ID))
}
