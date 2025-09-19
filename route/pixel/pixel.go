// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pixel

import (
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/wuhan005/NekoBox/internal/conf"
	"github.com/wuhan005/NekoBox/internal/context"
)

func Index(ctx context.Context) {
	ctx.Success("pixel")
}

func Proxy(ctx context.Context) error {
	uri := ctx.Param("**")
	method := ctx.Request().Method
	userID := strconv.Itoa(int(ctx.User.ID))

	var body io.Reader
	if method == http.MethodPost || method == http.MethodPut {
		body = ctx.Request().Request.Body
	}

	req, err := http.NewRequest(method, "http://pixel/", body)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create request")
		return ctx.ServerError()
	}
	req.URL.Host = conf.Pixel.Host
	req.URL.Path = path.Join("/api/", uri)
	req.Header.Set("neko-user-id", userID)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to send request")
		return ctx.ServerError()
	}
	defer func() { _ = resp.Body.Close() }()

	for k, v := range resp.Header {
		ctx.ResponseWriter().Header()[k] = v
	}
	ctx.ResponseWriter().WriteHeader(resp.StatusCode)

	_, err = io.Copy(ctx.ResponseWriter(), resp.Body)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to copy response")
		return ctx.ServerError()
	}

	return nil
}
