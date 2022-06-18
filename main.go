package main

import (
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

	db, err := db.Init()
	if err != nil {
		log.Fatal("Failed to connect database: %v", err)
	}

	r := route.New(db)

	r.Run("0.0.0.0", conf.Server.Port)
}
