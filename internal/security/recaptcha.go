// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/NekoWheel/NekoBox/internal/conf"
)

type recaptchaResponse struct {
	Success bool `json:"success"`
}

func CheckRecaptcha(response string, remoteIP string) (bool, error) {
	form := url.Values{}
	form.Set("secret", conf.Recaptcha.ServerKey)
	form.Set("response", response)
	form.Set("remoteip", remoteIP)
	requestBody := form.Encode()

	resp, err := http.Post(conf.Recaptcha.Domain+"/recaptcha/api/siteverify", "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	if err != nil {
		return false, errors.Wrap(err, "post request")
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var recaptchaResponse recaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&recaptchaResponse); err != nil {
		return false, errors.Wrap(err, "decode JSON")
	}

	return recaptchaResponse.Success, nil
}
