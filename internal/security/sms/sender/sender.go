// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sender

import (
	"context"

	"github.com/NekoWheel/NekoBox/internal/security/sms/storage"
)

type Sender interface {
	Send(ctx context.Context, typ storage.SMSType, phone, code string) error
}
