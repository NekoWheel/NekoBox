// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
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
		if avatarFileHeader.Size > storage.MaxAvatarSize {
			ctx.SetError(errors.New("头像文件太大，最大支持 2MB"))
			ctx.Success("user/profile")
			return
		}
		avatarURL, err = storage.UploadPicture(avatarFile, avatarFileHeader)
		if err != nil {
			log.Error("Failed to upload avatar: %v", err)
		}
	}

	var backgroundURL string
	backgroundFile, backgroundFileHeader, err := ctx.Request().FormFile("background")
	if err == nil {
		if backgroundFileHeader.Size > storage.MaxBackgroundSize {
			ctx.SetError(errors.New("背景图文件太大，最大支持 2MB"))
			ctx.Success("user/profile")
			return
		}
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

func ExportProfile(ctx context.Context) {
	user, err := db.Users.GetByID(ctx.Request().Context(), ctx.User.ID)
	if err != nil {
		log.Error("Failed to get user: %v", err)

		ctx.SetError(errors.New("导出失败：获取用户信息失败"))
		ctx.Success("user/profile")
		return
	}

	questions, err := db.Questions.GetByUserID(ctx.Request().Context(), user.ID, false)
	if err != nil {
		log.Error("Failed to get questions: %v", err)

		ctx.SetError(errors.New("导出失败：获取问题信息失败"))
		ctx.Success("user/profile")
		return
	}

	f, err := createExportExcelFile(user, questions)
	if err != nil {
		log.Error("Failed to create excel file: %v", err)

		ctx.SetError(errors.New("导出失败：创建Excel文件失败"))
		ctx.Success("user/profile")
		return
	}

	fileName := fmt.Sprintf("NekoBox账号信息导出-%s-%s.xlsx", user.Domain, time.Now().Format("20060102150405"))
	ctx.ResponseWriter().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.ResponseWriter().Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(fileName))

	if err := f.Write(ctx.ResponseWriter()); err != nil {
		log.Error("Failed to write excel file: %v", err)

		ctx.SetError(errors.New("导出失败：写入Excel文件失败"))
		ctx.Success("user/profile")
		return
	}
}

func createXLSXStreamWriter(xlsx *excelize.File, sheet string, headers []string) (*excelize.StreamWriter, error) {
	xlsx.NewSheet(sheet)
	sw, err := xlsx.NewStreamWriter(sheet)
	if err != nil {
		return nil, errors.Wrap(err, "new stream writer")
	}

	cols := make([]interface{}, 0, len(headers))
	for _, c := range headers {
		cols = append(cols, c)
	}
	err = sw.SetRow("A1", cols)
	if err != nil {
		return nil, errors.Wrap(err, "set header row")
	}
	return sw, nil
}

func createExportExcelFile(user *db.User, questions []*db.Question) (*excelize.File, error) {
	f := excelize.NewFile()

	sw, err := createXLSXStreamWriter(f, "账号信息", nil)
	if err != nil {
		return nil, errors.Wrap(err, "create xlsx stream writer: 提问")
	}
	// Set personal information sheet.
	personalData := [][]interface{}{
		{"NekoBox 账号信息导出", fmt.Sprintf("导出时间 %s", time.Now().Format("2006-01-02 15:04:05"))},
		{"电子邮箱", user.Email},
		{"昵称", user.Name},
		{"个性域名", user.Domain},
		{"介绍", user.Intro},
		{"头像 URL", user.Avatar},
		{"背景图 URL", user.Background},
		{"注册时间", user.CreatedAt},
	}
	currentRow := 1
	for _, row := range personalData {
		cell, _ := excelize.CoordinatesToCellName(1, currentRow)
		_ = sw.SetRow(cell, row)
		currentRow++
	}
	if err := sw.Flush(); err != nil {
		return nil, errors.Wrap(err, "flush personal data")
	}

	// Set questions sheet.
	sw, err = createXLSXStreamWriter(f, "提问", []string{"提问时间", "问题", "回答", "操作 Token"})
	if err != nil {
		return nil, errors.Wrap(err, "create xlsx stream writer: 提问")
	}

	currentRow = 2 // Include header row.
	for _, question := range questions {
		question := question
		vals := []interface{}{question.CreatedAt, question.Content, question.Answer, question.Token}
		cell, _ := excelize.CoordinatesToCellName(1, currentRow)
		_ = sw.SetRow(cell, vals)
		currentRow++
	}
	if err := sw.Flush(); err != nil {
		return nil, errors.Wrap(err, "flush personal data")
	}

	f.SetActiveSheet(f.GetSheetIndex("提问"))
	f.DeleteSheet("Sheet1") // Delete default sheet.

	return f, nil
}

func DeactivateProfile(ctx context.Context) {
	ctx.Success("user/deactivate")
}

func DeactivateProfileAction(ctx context.Context) {

}
