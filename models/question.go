package models

import (
	"errors"
)

func NewQuestion(form *QuestionForm) error {
	tx := DB.Begin()
	if tx.Create(&Question{
		PageID:  form.PageID,
		Content: form.Content,
		Answer:  "",
	}).RowsAffected != 1 {
		tx.Rollback()
		return errors.New("服务器错误！")
	}
	tx.Commit()
	return nil
}
