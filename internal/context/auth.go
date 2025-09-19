// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"github.com/flamego/flamego"
	"github.com/flamego/session"

	"github.com/wuhan005/NekoBox/internal/db"
)

// authenticatedUser returns the user object of the authenticated user.
func authenticatedUser(ctx flamego.Context, sess session.Session) *db.User {
	uid, ok := sess.Get("uid").(uint)
	if !ok {
		return nil
	}

	user, _ := db.Users.GetByID(ctx.Request().Context(), uid)
	return user
}

type ToggleOptions struct {
	UserSignInRequired  bool
	UserSignOutRequired bool
}

func Toggle(options *ToggleOptions) flamego.Handler {
	return func(ctx Context, endpoint EndpointType) error {
		if options.UserSignOutRequired && ctx.IsLogged {
			ctx.Redirect("/")
			return nil
		}

		if options.UserSignInRequired && !ctx.IsLogged {
			if endpoint.IsAPI() {
				return ctx.JSONError(40100, "请先登录")
			}
			ctx.SetErrorFlash("请先登录！")
			ctx.Redirect("/login")
			return nil
		}
		return nil
	}
}
