// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCensorLogs(t *testing.T) {
	t.Parallel()

	db, cleanup := newTestDB(t)
	ctx := context.Background()

	censorLogsStore := NewCensorLogsStore(db)

	for _, tc := range []struct {
		name string
		test func(t *testing.T, ctx context.Context, db *censorLogs)
	}{
		{"Create", testCensorLogsCreate},
		{"GetByText", testCensorLogsGetByText},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup("censor_logs"); err != nil {
					t.Fatal(err)
				}
			})
			tc.test(t, ctx, censorLogsStore.(*censorLogs))
		})
	}
}

func testCensorLogsCreate(t *testing.T, ctx context.Context, db *censorLogs) {
	err := db.Create(ctx, CreateCensorLogOptions{
		SourceName:  "aliyun",
		Input:       "hello world",
		Pass:        true,
		RawResponse: []byte(`{}`),
	})
	require.Nil(t, err)
}

func testCensorLogsGetByText(t *testing.T, ctx context.Context, db *censorLogs) {
	err := db.Create(ctx, CreateCensorLogOptions{
		SourceName:  "aliyun",
		Input:       "hello world",
		Pass:        true,
		RawResponse: []byte(`{}`),
	})
	require.Nil(t, err)

	t.Run("normal by cache", func(t *testing.T) {
		got, err := db.GetByText(ctx, "aliyun", "hello world")
		require.Nil(t, err)

		got.CreatedAt = time.Time{}
		got.UpdatedAt = time.Time{}

		want := &CensorLog{
			Model: gorm.Model{
				ID: 1,
			},
			SourceName:  "aliyun",
			Input:       "hello world",
			InputHash:   hashText("hello world"),
			Pass:        true,
			RawResponse: []byte(`{}`),
		}
		require.Equal(t, want, got)
	})

	t.Run("no longer than", func(t *testing.T) {
		_, err := db.GetByText(ctx, "aliyun", "hello world", time.Now())
		require.Equal(t, ErrCensorLogsNotFound, err)
	})
}
