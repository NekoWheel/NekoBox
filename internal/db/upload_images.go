// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var UploadImages UploadImagesStore

var _ UploadImagesStore = (*uploadImages)(nil)

type UploadImagesStore interface {
	Create(ctx context.Context, opts CreateUploadImageOptions) (*UploadImage, error)
	GetByQuestionID(ctx context.Context, questionID uint) ([]*UploadImage, error)
	GetByTypeQuestionID(ctx context.Context, typ UploadImageQuestionType, questionID uint) ([]*UploadImage, error)
}

func NewUploadImagesStore(db *gorm.DB) UploadImagesStore {
	return &uploadImages{db}
}

type uploadImages struct {
	*gorm.DB
}

type UploadImage struct {
	gorm.Model
	UploaderUserID uint
	Name           string
	FileSize       int64
	Md5            string
	Key            string
}

type UploadImageQuestionType string

const (
	UploadImageQuestionTypeAsk    UploadImageQuestionType = "ask"
	UploadImageQuestionTypeAnswer UploadImageQuestionType = "answer"
)

type UploadImageQuestion struct {
	Type          UploadImageQuestionType
	UploadImageID uint
	QuestionID    uint
}

type CreateUploadImageOptions struct {
	Type               UploadImageQuestionType
	QuestionID         uint
	UploaderUserID     uint
	Name               string
	FileSize           int64
	Md5                string
	Key                string
	IsDeletingPrevious bool
}

func (db *uploadImages) Create(ctx context.Context, opts CreateUploadImageOptions) (*UploadImage, error) {
	image := &UploadImage{
		UploaderUserID: opts.UploaderUserID,
		Name:           opts.Name,
		FileSize:       opts.FileSize,
		Md5:            opts.Md5,
		Key:            opts.Key,
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(image).Error; err != nil {
			return errors.Wrap(err, "create image")
		}

		if opts.IsDeletingPrevious {
			if err := tx.WithContext(ctx).Where("type = ? AND question_id = ?", opts.Type, opts.QuestionID).Delete(&UploadImageQuestion{}).Error; err != nil {
				return errors.Wrap(err, "delete previous image question link")
			}
		}
		if err := tx.WithContext(ctx).Create(&UploadImageQuestion{
			Type:          opts.Type,
			UploadImageID: image.ID,
			QuestionID:    opts.QuestionID,
		}).Error; err != nil {
			return errors.Wrap(err, "create image question link")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return image, nil
}

func (db *uploadImages) getBy(ctx context.Context, where string, args ...interface{}) ([]*UploadImage, error) {
	var uploadImageQuestions []*UploadImageQuestion
	if err := db.WithContext(ctx).Model(&UploadImageQuestion{}).Where(where, args...).Find(&uploadImageQuestions).Error; err != nil {
		return nil, errors.Wrap(err, "get image questions")
	}

	uploadImageIDs := lo.Map(uploadImageQuestions, func(item *UploadImageQuestion, _ int) uint {
		return item.UploadImageID
	})

	var uploadImages []*UploadImage
	if err := db.WithContext(ctx).Model(&UploadImage{}).Where("id IN ?", uploadImageIDs).Order("id ASC").Find(&uploadImages).Error; err != nil {
		return nil, errors.Wrap(err, "get images")
	}
	return uploadImages, nil
}

func (db *uploadImages) GetByQuestionID(ctx context.Context, questionID uint) ([]*UploadImage, error) {
	return db.getBy(ctx, "question_id = ?", questionID)
}

func (db *uploadImages) GetByTypeQuestionID(ctx context.Context, typ UploadImageQuestionType, questionID uint) ([]*UploadImage, error) {
	return db.getBy(ctx, "type = ? AND question_id = ?", typ, questionID)
}
