// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
)

func Login(ctx context.Context) {
	ctx.Success("auth/login")
}

func LoginAction(ctx context.Context, f form.Login, recaptcha recaptcha.RecaptchaV2) {
	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		log.Error("Failed to check recaptcha: %v", err)
		ctx.SetErrorFlash("内部错误，请稍后再试")
		ctx.Redirect("/login")
		return
	}
	if !resp.Success {
		ctx.SetErrorFlash("验证码错误")
		ctx.Redirect("/login")
		return
	}

	if ctx.HasError() {
		ctx.Success("auth/login")
		return
	}

	user, err := db.Users.Authenticate(ctx.Request().Context(), f.Email, f.Password)
	if err != nil {
		if errors.Is(err, db.ErrBadCredential) {
			ctx.SetErrorFlash(errors.Cause(err).Error())
		} else {
			log.Error("Failed to authenticate user: %v", err)
			ctx.SetErrorFlash("服务器错误！")
		}
		ctx.Redirect("/login")
		return
	}

	ctx.Session.Set("uid", user.ID)
	ctx.Redirect("/_/" + user.Domain)
}
