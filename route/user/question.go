// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
)

func QuestionList(ctx context.Context) {
	questions, err := db.Questions.GetByUserID(ctx.Request().Context(), ctx.User.ID, false)
	if err != nil {
		log.Error("Failed to get questions by user ID: %v", err)
		ctx.Redirect("/")
		return
	}
	ctx.Data["Questions"] = questions

	ctx.Success("user/question-list")
}
