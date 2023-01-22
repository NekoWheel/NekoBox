// Copyright 2021 E99p1ant. All rights reserved.

package dbutil

import (
	"time"
)

func Now() time.Time {
	return time.Now().Truncate(time.Microsecond)
}
