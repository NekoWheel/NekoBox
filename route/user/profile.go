// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"errors"

	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/storage"
)

func Profile(ctx context.Context) {
	ctx.Success("user/profile")
}

func UpdateProfile(ctx context.Context, f form.UpdateProfile) {
	if ctx.HasError() {
		ctx.Success("user/profile")
		return
	}

	var avatarURL string
	avatarFile, avatarFileHeader, err := ctx.Request().FormFile("avatar")
	if err == nil {
		avatarURL, err = storage.UploadPicture(avatarFile, avatarFileHeader)
		if err != nil {
			log.Error("Failed to upload avatar: %v", err)
		}
	}

	var backgroundURL string
	backgroundFile, backgroundFileHeader, err := ctx.Request().FormFile("background")
	if err == nil {
		backgroundURL, err = storage.UploadPicture(backgroundFile, backgroundFileHeader)
		if err != nil {
			log.Error("Failed to upload background: %v", err)
		}
	}

	if f.NewPassword != "" {
		if err := db.Users.ChangePassword(ctx.Request().Context(), ctx.User.ID, f.OldPassword, f.NewPassword); err != nil {
			if errors.Is(err, db.ErrBadCredential) {
				ctx.SetError(errors.New("旧密码输入错误"))
			} else {
				log.Error("Failed to update password: %v", err)
				ctx.SetError(errors.New("系统内部错误"))
			}
			ctx.Success("user/profile")
			return
		}
	}

	var notify db.NotifyType
	if f.NotifyEmail != "" {
		notify = db.NotifyTypeEmail
	} else {
		notify = db.NotifyTypeNone
	}

	if err := db.Users.Update(ctx.Request().Context(), ctx.User.ID, db.UpdateUserOptions{
		Name:       f.Name,
		Avatar:     avatarURL,
		Background: backgroundURL,
		Intro:      f.Intro,
		Notify:     notify,
	}); err != nil {
		ctx.SetErrorFlash("系统内部错误")
	} else {
		ctx.SetSuccessFlash("更新个人信息成功")
	}
	ctx.Redirect("/user/profile")
}
