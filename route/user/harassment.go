// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
)

func UpdateHarassment(ctx context.Context, f form.UpdateHarassment) {
	harassmentSetting := db.HarassmentSettingNone
	if f.RegisterOnly != "" {
		harassmentSetting = db.HarassmentSettingTypeRegisterOnly
	}

	if err := db.Users.UpdateHarassmentSetting(ctx.Request().Context(), ctx.User.ID, harassmentSetting); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to update harassment setting")
		ctx.SetInternalErrorFlash()
		ctx.Redirect("/user/profile")
		return
	}

	ctx.SetSuccessFlash("更新防骚扰设置成功")
	ctx.Redirect("/user/profile")
}
