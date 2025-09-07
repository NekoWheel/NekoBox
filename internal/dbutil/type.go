// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ContentCensorPass bool

func (ContentCensorPass) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "boolean GENERATED ALWAYS AS (IFNULL(content_censor_metadata->'$.pass' = true, false)) STORED NOT NULL"
	case "postgres":
		return "BOOLEAN GENERATED ALWAYS AS (COALESCE(content_censor_metadata->>'pass' = 'true', false)) STORED"
	}
	return ""
}

func (c *ContentCensorPass) Scan(value interface{}) error {
	var i sql.NullBool
	if err := i.Scan(value); err != nil {
		return err
	}
	*c = ContentCensorPass(i.Bool)
	return nil
}

type AnswerCensorPass bool

func (AnswerCensorPass) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "boolean GENERATED ALWAYS AS (IFNULL(answer_censor_metadata->'$.pass' = true, false)) STORED NOT NULL"
	case "postgres":
		return "BOOLEAN GENERATED ALWAYS AS (COALESCE(answer_censor_metadata->>'pass' = 'true', false)) STORED"
	}
	return ""
}

func (c *AnswerCensorPass) Scan(value interface{}) error {
	var i sql.NullBool
	if err := i.Scan(value); err != nil {
		return err
	}
	*c = AnswerCensorPass(i.Bool)
	return nil
}
