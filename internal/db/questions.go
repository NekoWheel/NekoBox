// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

var Questions QuestionsStore

type QuestionsStore interface {
	Create(ctx context.Context, opts CreateQuestionOptions) error
	GetByID(ctx context.Context, id uint) (*Question, error)
	GetByUserID(ctx context.Context, userID uint, answered bool) ([]*Question, error)
	AnswerByID(ctx context.Context, id uint, answer string) error
	DeleteByID(ctx context.Context, id uint) error
}

func NewQuestionsStore(db *gorm.DB) QuestionsStore {
	return &questions{db}
}

type questions struct {
	*gorm.DB
}

type Question struct {
	gorm.Model
	UserID  uint
	Content string
	Token   string
	Answer  string
}

type CreateQuestionOptions struct {
	UserID  uint
	Content string
}

func (db *questions) Create(ctx context.Context, opts CreateQuestionOptions) error {
	question := Question{
		UserID:  opts.UserID,
		Token:   randstr.String(6),
		Content: opts.Content,
	}
	return db.WithContext(ctx).Create(&question).Error
}

var ErrQuestionNotExist = errors.New("提问不存在")

func (db *questions) GetByID(ctx context.Context, id uint) (*Question, error) {
	var question Question
	if err := db.WithContext(ctx).First(&question, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotExist
		}
		return nil, errors.Wrap(err, "get question by ID")
	}
	return &question, nil
}

func (db *questions) GetByUserID(ctx context.Context, userID uint, answered bool) ([]*Question, error) {
	var questions []*Question
	q := db.WithContext(ctx)
	if answered {
		q = q.Where(`user_id = ? AND answer <> ""`, userID)
	} else {
		q = q.Where(`user_id = ?`, userID)
	}

	if err := q.Find(&questions).Error; err != nil {
		return nil, errors.Wrap(err, "get questions by page ID")
	}
	return questions, nil
}

func (db *questions) AnswerByID(ctx context.Context, id uint, answer string) error {
	var question Question
	if err := db.WithContext(ctx).First(&question, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotExist
		}
		return errors.Wrap(err, "get question by ID")
	}

	if err := db.WithContext(ctx).Model(&question).Where("id = ?", id).Update("answer", answer).Error; err != nil {
		return errors.Wrap(err, "update question answer")
	}
	return nil
}

func (db *questions) DeleteByID(ctx context.Context, id uint) error {
	var question Question
	if err := db.WithContext(ctx).First(&question, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotExist
		}
		return errors.Wrap(err, "get question by ID")
	}

	if err := db.WithContext(ctx).Delete(&Question{}, id).Error; err != nil {
		return errors.Wrap(err, "delete question")
	}
	return nil
}
