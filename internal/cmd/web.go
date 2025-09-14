// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"github.com/urfave/cli/v2"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/route"
	"github.com/NekoWheel/NekoBox/internal/tracing"
)

var Web = &cli.Command{
	Name:   "web",
	Usage:  "Start web server",
	Action: runWeb,
}

func runWeb(ctx *cli.Context) error {
	if err := conf.Init(); err != nil {
		return errors.Wrap(err, "load configuration")
	}

	if conf.Tracing.Enabled {
		logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))

		otelShutdown, err := tracing.SetupOTelSDK(ctx.Context, tracing.SetupOTelSDKOptions{
			TracingEndpoint: conf.Tracing.Endpoint,
			TracingToken:    conf.Tracing.Token,
			ServiceName:     conf.Tracing.ServiceName,
			HostName:        os.Getenv("HOSTNAME"),
		})
		if err != nil {
			logrus.WithContext(ctx.Context).WithError(err).Fatal("Failed to initialize OTel SDK")
		}

		defer func() {
			_ = otelShutdown(context.Background())
		}()
		logrus.WithContext(ctx.Context).Debug("Tracing enabled.")
	}

	dbType := conf.Database.Type

	var dsn string
	switch dbType {
	case "mysql", "":
		dsn = conf.MySQLDsn()
	case "postgres":
		dsn = conf.PostgresDsn()
	default:
		return errors.Errorf("unknown database type: %q", dbType)
	}
	conf.Database.DSN = dsn

	_, err := db.Init(dbType, dsn)
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}

	logrus.WithContext(ctx.Context).WithField("external_url", conf.App.ExternalURL).Info("Starting web server")
	r := route.New()
	r.Use(tracing.Middleware("NekoBox"))
	r.Run(conf.Server.Port)

	return nil
}
