// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var flagParseOnce sync.Once

func NewTestDB(t *testing.T, migrationTables ...interface{}) (testDB *gorm.DB, cleanup func(...string) error) {
	dsn := os.ExpandEnv("$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_DATABASE?charset=utf8mb4&parseTime=True&loc=Local")
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc:                Now,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to open connection: %v", err)
	}

	ctx := context.Background()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dbname := "cardinal-test-" + strconv.FormatUint(rng.Uint64(), 10)

	err = db.WithContext(ctx).Exec(`CREATE DATABASE ` + QuoteIdentifier(dbname)).Error
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// HACK: replace the database name in the DSN.
	dsn = strings.ReplaceAll(dsn, os.Getenv("DB_DATABASE"), dbname)

	flagParseOnce.Do(flag.Parse)

	testDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc:                Now,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to open test connection: %v", err)
	}

	err = testDB.AutoMigrate(migrationTables...)
	if err != nil {
		t.Fatalf("Failed to auto migrate tables: %v", err)
	}

	t.Cleanup(func() {
		defer func() {
			if database, err := db.DB(); err == nil {
				_ = database.Close()
			}
		}()

		if t.Failed() {
			t.Logf("DATABASE %s left intact for inspection", dbname)
			return
		}

		database, err := testDB.DB()
		if err != nil {
			t.Fatalf("Failed to get currently open database: %v", err)
		}

		err = database.Close()
		if err != nil {
			t.Fatalf("Failed to close currently open database: %v", err)
		}

		err = db.WithContext(ctx).Exec(`DROP DATABASE ` + QuoteIdentifier(dbname)).Error
		if err != nil {
			t.Fatalf("Failed to drop test database: %v", err)
		}
	})

	return testDB, func(tables ...string) error {
		if t.Failed() {
			return nil
		}

		for _, table := range tables {
			err := testDB.WithContext(ctx).Exec(`TRUNCATE TABLE ` + QuoteIdentifier(table)).Error
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// QuoteIdentifier quotes an "identifier" (e.g. a table or a column name) to be
// used as part of an SQL statement.
func QuoteIdentifier(s string) string {
	return "`" + strings.ReplaceAll(s, "`", "``") + "`"
}
