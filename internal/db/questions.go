// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var Questions QuestionsStore

type QuestionsStore interface {
	Create(ctx context.Context, opts CreateQuestionOptions) (*Question, error)
	GetByID(ctx context.Context, id uint) (*Question, error)
	GetByUserID(ctx context.Context, userID uint, answered bool) ([]*Question, error)
	AnswerByID(ctx context.Context, id uint, answer string) error
	DeleteByID(ctx context.Context, id uint) error
	UpdateCensor(ctx context.Context, id uint, opts UpdateQuestionCensorOptions) error
}

func NewQuestionsStore(db *gorm.DB) QuestionsStore {
	return &questions{db}
}

type questions struct {
	*gorm.DB
}

type Question struct {
	gorm.Model
	FromIP                string
	UserID                uint
	Content               string
	ContentCensorMetadata datatypes.JSON
	ContentCensorPass     bool `gorm:"->;type:boolean GENERATED ALWAYS AS (IFNULL(content_censor_metadata->'$.pass' = true, false)) STORED NOT NULL"`
	Token                 string
	Answer                string
	AnswerCensorMetadata  datatypes.JSON
	AnswerCensorPass      bool `gorm:"->;type:boolean GENERATED ALWAYS AS (IFNULL(answer_censor_metadata->'$.pass' = true, false)) STORED NOT NULL"`
}

type CreateQuestionOptions struct {
	FromIP  string
	UserID  uint
	Content string
}

func (db *questions) Create(ctx context.Context, opts CreateQuestionOptions) (*Question, error) {
	question := Question{
		FromIP:  opts.FromIP,
		UserID:  opts.UserID,
		Token:   randstr.String(6),
		Content: opts.Content,
	}
	return &question, db.WithContext(ctx).Create(&question).Error
}

type UpdateQuestionCensorOptions struct {
	ContentCensorMetadata json.RawMessage
	AnswerCensorMetadata  json.RawMessage
}

func (db *questions) UpdateCensor(ctx context.Context, id uint, opts UpdateQuestionCensorOptions) error {
	question, err := db.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get by ID")
	}

	contentCensorMetadata := question.ContentCensorMetadata
	if checkTextCensorResponseValid(opts.ContentCensorMetadata) {
		contentCensorMetadata = datatypes.JSON(opts.ContentCensorMetadata)
	}
	answerCensorMetadata := question.AnswerCensorMetadata
	if checkTextCensorResponseValid(opts.AnswerCensorMetadata) {
		answerCensorMetadata = datatypes.JSON(opts.AnswerCensorMetadata)
	}

	return db.WithContext(ctx).Model(&Question{}).Where("id = ?", id).Updates(&Question{
		ContentCensorMetadata: contentCensorMetadata,
		AnswerCensorMetadata:  answerCensorMetadata,
	}).Error
}

func checkTextCensorResponseValid(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}

	if bytes.EqualFold(raw, []byte("null")) {
		return false
	}

	var response struct {
		SourceName string `json:"source_name"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		return false
	}
	return response.SourceName != ""
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

	if err := q.Order("created_at DESC").Find(&questions).Error; err != nil {
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
