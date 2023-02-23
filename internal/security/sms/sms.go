// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sms

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/security/sms/sender"
	"github.com/NekoWheel/NekoBox/internal/security/sms/storage"
)

type SMSer interface {
	Send(ctx context.Context, phone string) error
	Validate(ctx context.Context, phone, code string) (bool, error)
}

type SMS struct {
	sender  sender.Sender
	storage storage.Storage
}

func NewSMS(sender sender.Sender, storage storage.Storage) *SMS {
	return &SMS{
		sender:  sender,
		storage: storage,
	}
}

func (s *SMS) Send(ctx context.Context, typ storage.SMSType, phone string) error {
	code := GenerateCode()
	logrus.WithContext(ctx).WithField("phone", phone).WithField("code", code).Info("send phone validation code")

	if err := s.sender.Send(ctx, typ, phone, code); err != nil {
		return errors.Wrap(err, "send")
	}
	if err := s.storage.Create(ctx, typ, phone, code); err != nil {
		return errors.Wrap(err, "create storage")
	}
	return nil
}

func (s *SMS) Validate(ctx context.Context, typ storage.SMSType, phone, code string) (bool, error) {
	return s.storage.Validate(ctx, typ, phone, code)
}

// GenerateCode generates a random phone validation code.
func GenerateCode() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", random.Int31n(1000000))
}
