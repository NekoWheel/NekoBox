package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"github.com/uptrace/opentelemetry-go-extra/otelplay"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/route"
)

func main() {
	logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	)))

	if err := conf.Init(); err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	_, err := db.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect database")
	}

	ctx := context.Background()
	shutdown := otelplay.ConfigureOpentelemetry(ctx)
	defer shutdown()

	r := route.New()

	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.App.UptraceDSN),
		uptrace.WithServiceName("neko-box-http-server"),
		uptrace.WithServiceVersion(conf.BuildCommit),
	)
	handler := otelhttp.NewHandler(r, "NekoBox")
	server := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(conf.Server.Port),
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		logrus.WithContext(ctx).WithError(err).Fatal("Failed to start server")
	}
}
