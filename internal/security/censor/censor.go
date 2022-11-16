// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package censor

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
)

var (
	censorCacheNoMoreThan = 31 * 24 * time.Hour // 1 month
)

// Text checks the text for sensitive content.
// It will save the censor log to the database and invoke the callback function.
func Text(ctx context.Context, text string) (*TextCensorResponse, error) {
	// Try to get the censor from the database log.
	censorLog, err := db.CensorLogs.GetByText(ctx, text, time.Now().Add(-censorCacheNoMoreThan))
	if err != nil {
		if !errors.Is(err, db.ErrCensorLogsNotFound) {
			logrus.WithContext(ctx).WithError(err).Error("Failed to get censor log")
		}
	} else {
		raw := censorLog.RawResponse
		// We got the previous censor log cache, seems like we saved money for the API call.
		switch censorLog.SourceName {
		case "qiniu":
			return QiniuTextCensorParser(raw)
		default:
			logrus.WithContext(ctx).WithField("censor_source", censorLog.SourceName).Error("Unknown censor source")
		}
	}

	// TODO: support more censor providers selection.
	sourceName := "qiniu"
	censor := NewQiniuTextCensor(conf.App.QiniuAccessKey, conf.App.QiniuAccessSecret)
	response, err := censor.Censor(ctx, text)
	if err != nil {
		return nil, errors.Wrap(err, "censor")
	}

	// Save the censor log to the database.
	if err := db.CensorLogs.Create(ctx, db.CreateCensorLogOptions{
		SourceName:  sourceName,
		Input:       text,
		Pass:        response.Pass,
		RawResponse: response.RawResponse,
	}); err != nil {
		logrus.WithContext(ctx).WithError(err).Error("Failed to save censor log")
	}

	return response, nil
}
