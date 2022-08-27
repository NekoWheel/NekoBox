// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"

	"github.com/NekoWheel/NekoBox/internal/conf"
)

const (
	OSSPictureKeyPrefix = "picture/"
)

// UploadPictureToOSS upload user's avatar or background image to OSS.
// It returns the uploaded asset URL.
func UploadPictureToOSS(file multipart.File, _ *multipart.FileHeader) (string, error) {
	client, err := oss.New(conf.Upload.AliyunEndpoint, conf.Upload.AliyunAccessID, conf.Upload.AliyunAccessSecret)
	if err != nil {
		return "", errors.Wrap(err, "new oss client")
	}

	bucket, err := client.Bucket(conf.Upload.AliyunBucket)
	if err != nil {
		return "", errors.Wrap(err, "bucket")
	}

	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()

	key := fmt.Sprintf("%s%d/%02d/%02d/%s", OSSPictureKeyPrefix, year, month, day, randstr.Hex(15))
	if err := bucket.PutObject(key, file); err != nil {
		return "", errors.Wrap(err, "put object")
	}

	return fmt.Sprintf("https://%s.%s/%s", conf.Upload.AliyunBucket, conf.Upload.AliyunEndpoint, key), nil
}
