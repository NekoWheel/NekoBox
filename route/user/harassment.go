// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
)

func UpdateHarassment(ctx context.Context, f form.UpdateHarassment) {
	harassmentSettingType := db.HarassmentSettingNone
	if f.RegisterOnly != "" {
		harassmentSettingType = db.HarassmentSettingTypeRegisterOnly
	}

	blockWords := f.BlockWords
	blockWords = strings.ReplaceAll(blockWords, "，", ",")
	blockWords = strings.TrimSpace(blockWords)

	words := make([]string, 0)
	wordSet := make(map[string]struct{})
	for _, word := range strings.Split(blockWords, ",") {
		word := strings.TrimSpace(word)
		if word == "" {
			continue
		}
		if _, ok := wordSet[word]; ok {
			continue
		}
		wordSet[word] = struct{}{}

		if len(word) > 10 {
			ctx.SetErrorFlash(fmt.Sprintf("屏蔽词长度不能超过 10 个字符：%s", word))
			ctx.Redirect("/user/profile")
			return
		}
		words = append(words, word)
	}

	if len(words) > 10 {
		ctx.SetErrorFlash("屏蔽词不能超过 10 个")
		ctx.Redirect("/user/profile")
		return
	}

	if err := db.Users.UpdateHarassmentSetting(ctx.Request().Context(), ctx.User.ID, db.HarassmentSettingOptions{
		Type:       harassmentSettingType,
		BlockWords: strings.Join(words, ","),
	}); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to update harassment setting")
		ctx.SetInternalErrorFlash()
		ctx.Redirect("/user/profile")
		return
	}

	ctx.SetSuccessFlash("更新防骚扰设置成功")
	ctx.Redirect("/user/profile")
}
