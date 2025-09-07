package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/urfave/cli/v2"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
)

var Uid = &cli.Command{
	Name:   "uid",
	Usage:  "Create uid for user",
	Action: runUid,
}

func runUid(ctx *cli.Context) error {
	if err := conf.Init(); err != nil {
		return errors.Wrap(err, "load configuration")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)
	dbType := conf.Database.Type
	conf.Database.DSN = dsn

	database, err := db.Init(dbType, dsn)
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}

	var users []*db.User
	if err := database.WithContext(ctx.Context).Unscoped().Find(&users).Error; err != nil {
		return errors.Wrap(err, "query users")
	}

	for idx, user := range users {
		user := user
		if user.UID != "" {
			continue
		}

		if idx%1000 == 0 {
			fmt.Printf("Processing user %d/%d\n", idx, len(users))
		}

		uid := xid.NewWithTime(user.CreatedAt).String()
		if err := database.WithContext(ctx.Context).Unscoped().
			Model(&db.User{}).
			Omit("updated_at").
			Where("id = ?", user.ID).UpdateColumn("uid", uid).Error; err != nil {
			return errors.Wrap(err, "update uid")
		}
	}
	return nil
}
