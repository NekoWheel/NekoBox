// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"os"

	"github.com/flamego/cache"
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/security/sms"
	"github.com/NekoWheel/NekoBox/route"
)

func Register(ctx context.Context) {
	ctx.Success("auth/register")
}

func SendRegisterSMSAPI(ctx context.Context, f form.SendSMS, sms sms.SMS, cache cache.Cache, recaptcha recaptcha.RecaptchaV3) error {
	return route.SendSMS(route.SMSCacheKeyPrefixRegister)(ctx, f, sms, cache, recaptcha)
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

	var phone string
	verifyType := db.VerifyTypeUnverified
	// Check sms code.
	if f.Phone != "" && f.VerifyCode != "" {
		verifyCodeInf, err := cache.Get(ctx.Request().Context(), route.SMSCacheKeyPrefixRegister+f.Phone)
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
				// Remove the key.
				if err := cache.Delete(ctx.Request().Context(), route.SMSCacheKeyPrefixRegister+f.Phone); err != nil {
					logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to delete register code cache")
				}

				// Set user's verified phone.
				phone = f.Phone
				verifyType = db.VerifyTypeVerified
			}
		}
	}

	if err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		Name:       f.Name,
		Password:   f.Password,
		Email:      f.Email,
		Phone:      phone,
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
			errors.Is(err, db.ErrDuplicateDomain),
			errors.Is(err, db.ErrDuplicatePhone):
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
