package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/uptrace/opentelemetry-go-extra/otelplay"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	log "unknwon.dev/clog/v2"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/route"
)

func main() {
	defer log.Stop()
	if err := log.NewConsole(); err != nil {
		panic("init console logger: " + err.Error())
	}

	if err := conf.Init(); err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	_, err := db.Init()
	if err != nil {
		log.Fatal("Failed to connect database: %v", err)
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
		log.Fatal("Failed to start server: %v", err)
	}
}
