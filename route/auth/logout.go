// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/NekoWheel/NekoBox/internal/context"
)

func Logout(ctx context.Context) {
	ctx.Session.Flush()
	ctx.Redirect("/")
}
