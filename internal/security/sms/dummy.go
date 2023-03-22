package sms

import (
	"context"

	"github.com/sirupsen/logrus"
)

type DummySMS struct{}

func NewDummySMS() *DummySMS {
	return &DummySMS{}
}

func (s *DummySMS) SendCode(ctx context.Context, phone, code string) error {
	logrus.WithContext(ctx).WithField("phone", phone).WithField("code", code).Trace("Send code to phone number, but do nothing")
	return nil
}
