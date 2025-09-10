// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"embed"
)

//go:embed auth base mail question user home.html sponsor.html change-logs.html pixel.html maintenance-mode.html
var FS embed.FS
