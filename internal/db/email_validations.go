// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"gorm.io/gorm"
)

type EmailValidation struct {
	gorm.Model

	UserID uint
	Email  string
	Code   string
	Type   string
}
