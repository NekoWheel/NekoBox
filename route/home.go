// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/wuhan005/NekoBox/internal/context"
)

func Home(ctx context.Context) {
	ctx.Success("home")
}

func Sponsor(ctx context.Context) {
	ctx.Success("sponsor")
}

func ChangeLogs(ctx context.Context) {
	ctx.Success("change-logs")
}
