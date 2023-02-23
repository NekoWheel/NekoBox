// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/security/sms"
	"github.com/NekoWheel/NekoBox/internal/security/sms/storage"
)

func Register(ctx context.Context) {
	ctx.Success("auth/register")
}

func RegisterAction(ctx context.Context, f form.Register, recaptcha recaptcha.RecaptchaV3, sms *sms.SMS) {
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

	// Check the phone SMS code.
	ok, err := sms.Validate(ctx.Request().Context(), storage.SMSTypeVerifyPhone, f.Phone, f.SMSCode)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to validate phone SMS code")
		ctx.SetInternalErrorFlash()
		ctx.Redirect("/register")
		return
	}
	if !ok {
		ctx.SetErrorFlash("手机验证码错误")
		ctx.Redirect("/register")
		return
	}

	if err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		Name:       f.Name,
		Password:   f.Password,
		Phone:      f.Phone,
		Email:      f.Email,
		Avatar:     conf.Upload.DefaultAvatarURL,
		Domain:     f.Domain,
		Background: conf.Upload.DefaultBackground,
		Intro:      "问你想问的",
	}); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotExists),
			errors.Is(err, db.ErrBadCredential),
			errors.Is(err, db.ErrDuplicatePhone),
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

func SendRegisterSMSCodeAPI(ctx context.Context, f form.SendSMS, sms *sms.SMS) error {
	if err := sms.Send(ctx.Request().Context(), storage.SMSTypeVerifyPhone, f.Phone); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to send SMS code")
		return ctx.ServerError()
	}
	return ctx.JSON("ok")
}
