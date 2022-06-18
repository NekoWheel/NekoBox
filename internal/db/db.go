// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/NekoWheel/NekoBox/internal/conf"
)

func Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Address,
		conf.Database.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "connect to database")
	}

	if err := db.AutoMigrate(&User{}, &Question{}); err != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	Users = NewUsersStore(db)
	Questions = NewQuestionsStore(db)

	return db, nil
}