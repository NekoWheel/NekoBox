// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Storage = (*Postgres)(nil)

type Postgres struct {
	db *gorm.DB
}

type SMS struct {
	gorm.Model
	Type       SMSType
	Phone      string
	Code       string `gorm:"index"`
	ExpiredAt  time.Time
	VerifiedAt time.Time
}

func (s *SMS) IsValid() bool {
	return s.VerifiedAt.IsZero() && s.ExpiredAt.After(time.Now())
}

func NewPostgresStorage(db *gorm.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (s *Postgres) Create(ctx context.Context, typ SMSType, phone, code string) error {
	return s.db.WithContext(ctx).Model(&SMS{}).Create(&SMS{
		Type:      typ,
		Phone:     phone,
		Code:      code,
		ExpiredAt: time.Now().Add(ExpiredTime),
	}).Error
}

func (s *Postgres) Validate(ctx context.Context, typ SMSType, phone, code string) (bool, error) {
	var sms SMS
	if err := s.db.WithContext(ctx).Model(&SMS{}).Where("type = ? AND phone = ? AND code = ?", typ, phone, code).First(&sms).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	// Check the code is expired.
	return sms.IsValid(), nil
}
