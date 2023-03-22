// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/flamego/cache"
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/security/sms"
	"github.com/NekoWheel/NekoBox/internal/strutil"
)

func Register(ctx context.Context) {
	ctx.Success("auth/register")
}

const (
	smsCodeCacheKeyPrefix = "register-sms-code:"
)

func SendRegisterSMS(ctx context.Context, f form.RegisterSendSMS, sms sms.SMS, cache cache.Cache, recaptcha recaptcha.RecaptchaV3) {
	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
		_ = ctx.ServerError()
		return
	}
	if !resp.Success {
		_ = ctx.JSONError(http.StatusBadRequest, "验证码错误")
		return
	}

	phone := f.Phone
	code := strutil.RandomNumericString(6)
	smsCodeCacheKey := smsCodeCacheKeyPrefix + phone
	if err := cache.Set(ctx.Request().Context(), smsCodeCacheKey, code, 5*time.Minute); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to set sms code cache")
		_ = ctx.ServerError()
		return
	}

	if err := sms.SendCode(ctx.Request().Context(), phone, code); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to send sms code")
		_ = ctx.JSONError(http.StatusBadRequest, "发送短信验证码失败，请稍后重试")
		return
	}

	logrus.WithContext(ctx.Request().Context()).WithField("phone", phone).WithField("code", code).Info("Send sms code successfully")
	_ = ctx.JSON("发送短信验证码成功")
}

func RegisterAction(ctx context.Context, f form.Register, cache cache.Cache, recaptcha recaptcha.RecaptchaV3) {
	if ctx.HasError() {
		ctx.Success("auth/register")
		return
	}

	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
		ctx.SetInternalErrorFlash()
		ctx.Redirect("/register")
		return
	}
	if !resp.Success {
		ctx.SetErrorFlash("验证码错误")
		ctx.Redirect("/register")
		return
	}

	verifyType := db.VerifyTypeUnverified
	// Check sms code.
	if f.Phone != "" && f.VerifyCode != "" {
		verifyCodeInf, err := cache.Get(ctx.Request().Context(), smsCodeCacheKeyPrefix+f.Phone)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				ctx.SetErrorFlash("验证码已过期")
			} else {
				logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to read password recovery code cache")
				ctx.SetInternalErrorFlash()
			}
			ctx.Redirect("/register")
			return
		} else {
			verifyCode, ok := verifyCodeInf.(string)
			if ok && verifyCode != "" && verifyCode == f.VerifyCode {
				verifyType = db.VerifyTypeVerified
			}
		}
	}

	if err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		Name:       f.Name,
		Password:   f.Password,
		Email:      f.Email,
		Avatar:     conf.Upload.DefaultAvatarURL,
		Domain:     f.Domain,
		Background: conf.Upload.DefaultBackground,
		Intro:      "问你想问的",
		VerifyType: verifyType,
	}); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotExists),
			errors.Is(err, db.ErrBadCredential),
			errors.Is(err, db.ErrDuplicateEmail),
			errors.Is(err, db.ErrDuplicateDomain):
			ctx.SetError(errors.Cause(err))

		default:
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create new user")
			ctx.SetInternalError()
		}

		ctx.Success("auth/register")
		return
	}

	ctx.SetSuccessFlash("注册成功，欢迎来到 NekoBox！")
	ctx.Redirect("/login")
}
