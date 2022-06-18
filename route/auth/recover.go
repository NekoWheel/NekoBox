// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"os"
	"time"

	"github.com/flamego/cache"
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/mail"
)

func ForgotPassword(ctx context.Context) {
	ctx.Success("auth/forgot-password")
}

func ForgotPasswordAction(ctx context.Context, f form.ForgotPassword, cache cache.Cache, recaptcha recaptcha.RecaptchaV2) {
	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		log.Error("Failed to check recaptcha: %v", err)
		ctx.SetErrorFlash("内部错误，请稍后再试")
		ctx.Redirect("/forgot-password")
		return
	}
	if !resp.Success {
		ctx.SetErrorFlash("验证码错误")
		ctx.Redirect("/forgot-password")
		return
	}

	if ctx.HasError() {
		ctx.Success("auth/forgot-password")
		return
	}

	user, err := db.Users.GetByEmail(ctx.Request().Context(), f.Email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotExists) {
			ctx.SetErrorFlash("用户邮箱不存在")
		} else {
			ctx.SetErrorFlash("内部错误，请稍后再试")
		}
		ctx.Redirect("/forgot-password")
		return
	}

	emailSentCacheKey := "forgot-password-email-sent:" + user.Email
	_, err = cache.Get(ctx.Request().Context(), emailSentCacheKey)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Error("Failed to read password recovery email sent cache: %v", err)
		}
	} else {
		ctx.SetErrorFlash("邮件发送太频繁，请稍后再试")
		ctx.Redirect("/forgot-password")
		return
	}

	code := randstr.String(64)
	recoveryCodeCacheKey := "forgot-password-recovery-code:" + code
	if err := cache.Set(ctx.Request().Context(), recoveryCodeCacheKey, user.ID, 24*time.Hour); err != nil {
		log.Error("Failed to set password recovery code cache: %v", err)
		ctx.SetErrorFlash("内部错误，请稍后再试")
		ctx.Redirect("/forgot-password")
		return
	}

	if err := mail.SendPasswordRecoveryMail(user.Email, code); err != nil {
		log.Error("Failed to send password recovery mail: %v", err)
		ctx.SetErrorFlash("邮件发送失败，请稍后再试")
		ctx.Redirect("/forgot-password")
		return
	}

	if err := cache.Set(ctx.Request().Context(), emailSentCacheKey, time.Now(), 2*time.Minute); err != nil {
		log.Error("Failed to set password recovery email cache: %v", err)
	}

	ctx.Data["email"] = user.Email
	ctx.Success("auth/forgot-password-sent")
}

func checkRecoverPasswordCode(ctx context.Context, cache cache.Cache) (*db.User, bool) {
	code := ctx.Query("code")
	recoveryCodeCacheKey := "forgot-password-recovery-code:" + code
	userIDItf, err := cache.Get(ctx.Request().Context(), recoveryCodeCacheKey)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ctx.SetErrorFlash("验证码已过期")
		} else {
			log.Error("Failed to read password recovery code cache: %v", err)
			ctx.SetErrorFlash("内部错误，请稍后再试")
		}
		ctx.Redirect("/login")
		return nil, false
	}

	userID, ok := userIDItf.(uint)
	if !ok {
		ctx.SetErrorFlash("内部错误，请稍后再试")
		ctx.Redirect("/login")
		return nil, false
	}

	user, err := db.Users.GetByID(ctx.Request().Context(), userID)
	if err != nil {
		ctx.SetErrorFlash("用户不存在")
		ctx.Redirect("/login")
		return nil, false
	}

	return user, true
}

func RecoverPassword(ctx context.Context, cache cache.Cache) {
	user, ok := checkRecoverPasswordCode(ctx, cache)
	if !ok {
		return
	}

	ctx.Data["User"] = user
	ctx.Success("auth/password-recovery")
}

func RecoverPasswordAction(ctx context.Context, cache cache.Cache, f form.RecoverPassword) {
	user, ok := checkRecoverPasswordCode(ctx, cache)
	if !ok {
		return
	}
	if ctx.HasError() {
		ctx.Data["User"] = user
		ctx.Success("auth/password-recovery")
		return
	}

	code := ctx.Query("code")
	recoveryCodeCacheKey := "forgot-password-recovery-code:" + code
	if err := cache.Delete(ctx.Request().Context(), recoveryCodeCacheKey); err != nil {
		log.Error("Failed to delete password recovery code cache: %v", err)
		ctx.SetErrorFlash("内部错误，请稍后再试")
		ctx.Redirect("/login")
		return
	}

	if err := db.Users.UpdatePassword(ctx.Request().Context(), user.ID, f.NewPassword); err != nil {
		ctx.SetErrorFlash("内部错误，请稍后再试")
		ctx.Refresh()
		return
	}

	ctx.SetSuccessFlash("密码修改成功")
	ctx.Redirect("/login")
}
