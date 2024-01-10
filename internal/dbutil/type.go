// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ContentCensorPass bool

func (ContentCensorPass) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "boolean GENERATED ALWAYS AS (IFNULL(content_censor_metadata->'$.pass' = true, false)) STORED NOT NULL"
	case "postgres":
		return "BOOLEAN GENERATED ALWAYS AS (COALESCE(content_censor_metadata->>'$.pass' = 'true', false)) STORED"
	}
	return ""
}

type AnswerCensorPass bool

func (AnswerCensorPass) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "boolean GENERATED ALWAYS AS (IFNULL(answer_censor_metadata->'$.pass' = true, false)) STORED NOT NULL"
	case "postgres":
		return "BOOLEAN GENERATED ALWAYS AS (COALESCE(answer_censor_metadata->>'$.pass' = 'true', false)) STORED"
	}
	return ""
}
