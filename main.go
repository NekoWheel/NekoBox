package main

import (
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/route"
)

func main() {
	if err := conf.Init(); err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	_, err := db.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect database")
	}

	r := route.New()

	r.Run("0.0.0.0", conf.Server.Port)
}
