// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/NekoWheel/NekoBox/internal/conf"
)

// MaxAvatarSize is the max avatar size which is 2MB.
const MaxAvatarSize = 2 * 1024 * 1024

// MaxBackgroundSize is the max background size which is 2MB.
const MaxBackgroundSize = 2 * 1024 * 1024

type uploadPictureCallBack struct {
	Code int `json:"code"`
	Data struct {
		Md5      string `json:"md5"`
		Mime     string `json:"mime"`
		Name     string `json:"name"`
		Quota    string `json:"quota"`
		Sha1     string `json:"sha1"`
		Size     int    `json:"size"`
		URL      string `json:"url"`
		UseQuota string `json:"use_quota"`
	} `json:"data"`
	Msg  string `json:"msg"`
	Time int    `json:"time"`
}

func UploadPicture(file multipart.File, header *multipart.FileHeader) (string, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("image", header.Filename)
	if err != nil {
		return "", errors.Wrap(err, "create form file")
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", errors.Wrap(err, "io copy")
	}
	err = writer.Close()
	if err != nil {
		return "", errors.Wrap(err, "close writer")
	}

	req, err := http.NewRequest(http.MethodPost, conf.Upload.URL, requestBody)
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}
	req.Header.Set("token", conf.Upload.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "request")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var callback uploadPictureCallBack
	if err := json.NewDecoder(resp.Body).Decode(&callback); err != nil {
		return "", errors.Wrap(err, "decode JSON")
	}

	url := strings.Split(callback.Data.URL, "?")[0]
	return url, nil
}
