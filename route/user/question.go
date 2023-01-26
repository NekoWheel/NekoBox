// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
)

func QuestionList(ctx context.Context) {
	questions, err := db.Questions.GetByUserID(ctx.Request().Context(), ctx.User.ID, db.GetQuestionsByUserIDOptions{
		FilterAnswered: false,
		ShowPrivate:    true,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get questions by user ID")
		ctx.Redirect("/")
		return
	}
	ctx.Data["Questions"] = questions

	ctx.Success("user/question-list")
}
