// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUploadImages(t *testing.T) {
	t.Parallel()

	db, cleanup := newTestDB(t)
	ctx := context.Background()

	uploadImagesStore := NewUploadImagesStore(db)

	for _, tc := range []struct {
		name string
		test func(t *testing.T, ctx context.Context, db *uploadImages)
	}{
		{"Create", testUploadImagesCreate},
		{"BindUploadImageWithQuestion", testUploadImagesBindUploadImageWithQuestion},
		{"GetByQuestionID", testUploadImagesGetByQuestionID},
		{"GetByTypeQuestionID", testUploadImagesGetByTypeQuestionID},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup("upload_images"); err != nil {
					t.Fatal(err)
				}
			})
			tc.test(t, ctx, uploadImagesStore.(*uploadImages))
		})
	}
}

func testUploadImagesCreate(t *testing.T, ctx context.Context, db *uploadImages) {
	got, err := db.Create(ctx, CreateUploadImageOptions{
		UploaderUserID: 1,
		Name:           "test.png",
		FileSize:       12345,
		Md5:            "d41d8cd98f00b204e9800998ecf8427e",
		Key:            "uploads/test.png",
	})
	require.Nil(t, err)

	want := &UploadImage{
		Model:          got.Model,
		UploaderUserID: 1,
		Name:           "test.png",
		FileSize:       12345,
		Md5:            "d41d8cd98f00b204e9800998ecf8427e",
		Key:            "uploads/test.png",
	}
	require.Equal(t, want, got)
}

func testUploadImagesBindUploadImageWithQuestion(t *testing.T, ctx context.Context, db *uploadImages) {
	_, err := db.Create(ctx, CreateUploadImageOptions{
		UploaderUserID: 1,
		Name:           "test.png",
		FileSize:       12345,
		Md5:            "d41d8cd98f00b204e9800998ecf8427e",
		Key:            "uploads/test.png",
	})
	require.Nil(t, err)

	_, err = db.Create(ctx, CreateUploadImageOptions{
		UploaderUserID: 1,
		Name:           "test2.png",
		FileSize:       12345,
		Md5:            "d41d8cd98f00b204e9800998ecf8427e",
		Key:            "uploads/test2.png",
	})
	require.Nil(t, err)

	t.Run("normal", func(t *testing.T) {
		err := db.BindUploadImageWithQuestion(ctx, 1, UploadImageQuestionTypeAsk, 1)
		require.Nil(t, err)
	})

	t.Run("bind answer", func(t *testing.T) {
		err := db.BindUploadImageWithQuestion(ctx, 2, UploadImageQuestionTypeAnswer, 1)
		require.Nil(t, err)

		got, err := db.GetByTypeQuestionID(ctx, UploadImageQuestionTypeAnswer, 1)
		require.Nil(t, err)

		want := []*UploadImage{{
			Model:          got[0].Model,
			UploaderUserID: 1,
			Name:           "test2.png",
			FileSize:       12345,
			Md5:            "d41d8cd98f00b204e9800998ecf8427e",
			Key:            "uploads/test2.png",
		}}
		require.Equal(t, got, want)
	})

	t.Run("bind both", func(t *testing.T) {
		err = db.BindUploadImageWithQuestion(ctx, 1, UploadImageQuestionTypeAsk, 1)
		require.Nil(t, err)

		err := db.BindUploadImageWithQuestion(ctx, 2, UploadImageQuestionTypeAnswer, 1)
		require.Nil(t, err)

		got, err := db.GetByTypeQuestionID(ctx, UploadImageQuestionTypeAsk, 1)
		require.Nil(t, err)

		want := []*UploadImage{{
			Model:          got[0].Model,
			UploaderUserID: 1,
			Name:           "test.png",
			FileSize:       12345,
			Md5:            "d41d8cd98f00b204e9800998ecf8427e",
			Key:            "uploads/test.png",
		}}
		require.Equal(t, got, want)

		got, err = db.GetByTypeQuestionID(ctx, UploadImageQuestionTypeAnswer, 1)
		require.Nil(t, err)

		want = []*UploadImage{{
			Model:          got[0].Model,
			UploaderUserID: 1,
			Name:           "test2.png",
			FileSize:       12345,
			Md5:            "d41d8cd98f00b204e9800998ecf8427e",
			Key:            "uploads/test2.png",
		}}
		require.Equal(t, got, want)
	})
}

func testUploadImagesGetByQuestionID(t *testing.T, ctx context.Context, db *uploadImages) {
	_, err := db.Create(ctx, CreateUploadImageOptions{
		UploaderUserID: 1,
		Name:           "test.png",
		FileSize:       12345,
		Md5:            "d41d8cd98f00b204e9800998ecf8427e",
		Key:            "uploads/test.png",
	})
	require.Nil(t, err)

	err = db.BindUploadImageWithQuestion(ctx, 1, UploadImageQuestionTypeAsk, 1)
	require.Nil(t, err)

	t.Run("normal", func(t *testing.T) {
		got, err := db.GetByQuestionID(ctx, 1)
		require.Nil(t, err)

		want := []*UploadImage{{
			Model:          got[0].Model,
			UploaderUserID: 1,
			Name:           "test.png",
			FileSize:       12345,
			Md5:            "d41d8cd98f00b204e9800998ecf8427e",
			Key:            "uploads/test.png",
		}}
		require.Equal(t, got, want)
	})

	t.Run("not found", func(t *testing.T) {
		got, err := db.GetByQuestionID(ctx, 2)
		require.Nil(t, err)
		require.Equal(t, 0, len(got))
	})
}

func testUploadImagesGetByTypeQuestionID(t *testing.T, ctx context.Context, db *uploadImages) {
	_, err := db.Create(ctx, CreateUploadImageOptions{
		UploaderUserID: 1,
		Name:           "test.png",
		FileSize:       12345,
		Md5:            "d41d8cd98f00b204e9800998ecf8427e",
		Key:            "uploads/test.png",
	})
	require.Nil(t, err)

	err = db.BindUploadImageWithQuestion(ctx, 1, UploadImageQuestionTypeAsk, 1)
	require.Nil(t, err)

	t.Run("normal", func(t *testing.T) {
		got, err := db.GetByTypeQuestionID(ctx, UploadImageQuestionTypeAsk, 1)
		require.Nil(t, err)

		want := []*UploadImage{{
			Model:          got[0].Model,
			UploaderUserID: 1,
			Name:           "test.png",
			FileSize:       12345,
			Md5:            "d41d8cd98f00b204e9800998ecf8427e",
			Key:            "uploads/test.png",
		}}
		require.Equal(t, got, want)
	})

	t.Run("not found", func(t *testing.T) {
		got, err := db.GetByTypeQuestionID(ctx, UploadImageQuestionTypeAsk, 2)
		require.Nil(t, err)
		require.Equal(t, 0, len(got))
	})
}
