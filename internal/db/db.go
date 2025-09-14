// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/NekoWheel/NekoBox/internal/dbutil"
)

var AllTables = []interface{}{
	&User{}, &Question{}, &CensorLog{}, &UploadImage{}, &UploadImageQuestion{},
}

func Init(typ, dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch typ {
	case "mysql", "":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		return nil, errors.Errorf("unknown database type: %q", typ)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return dbutil.Now()
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             3 * time.Second,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "connect to database")
	}

	if err := db.Use(tracing.NewPlugin(
		tracing.WithAttributes(
			attribute.String("db.name", db.Name()),
		),
	)); err != nil {
		return nil, errors.Wrap(err, "register otelgorm plugin")
	}

	if err := db.AutoMigrate(AllTables...); err != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	Users = NewUsersStore(db)
	Questions = NewQuestionsStore(db)
	CensorLogs = NewCensorLogsStore(db)
	UploadImages = NewUploadImagesStore(db)

	return db, nil
}
