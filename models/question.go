package models

import (
	"errors"
	"github.com/jinzhu/gorm"
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

func GetQuestionsByPageID(pageID uint, order bool) []*Question {
	questions := make([]*Question, 0)

	query := DB.Model(&Question{}).Where(&Question{
		PageID: pageID,
	})
	if order {
		query = query.Order("`answer` <> \"\", `id`")
	}
	query.Find(&questions)
	return questions
}

func GetQuestionByDomainID(domain string, questionID uint) (*Question, error) {
	page := new(Page)
	DB.Model(&Page{}).Where(&Page{Domain: domain}).Find(&page)
	if page.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	question := new(Question)
	DB.Model(&Question{}).Where(&Question{Model: gorm.Model{ID: questionID}, PageID: page.ID}).Find(&question)
	if question.ID == 0 {
		return nil, errors.New("问题不存在")
	}
	return question, nil
}

func AnswerQuestion(questionID uint, question *Question) error {
	tx := DB.Begin()
	if tx.Model(&Question{}).Where(&Question{Model: gorm.Model{ID: questionID}}).Update(&question).RowsAffected != 1 {
		tx.Rollback()
		return errors.New("回答问题失败")
	}
	tx.Commit()
	return nil
}
