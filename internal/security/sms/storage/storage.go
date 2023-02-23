// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"context"
	"time"
)

type SMSType string

const (
	ExpiredTime = 5 * time.Minute

	SMSTypeVerifyPhone SMSType = "verify_phone"
)

type Storage interface {
	Create(ctx context.Context, typ SMSType, phone, code string) error
	Validate(ctx context.Context, typ SMSType, phone, code string) (bool, error)
}
