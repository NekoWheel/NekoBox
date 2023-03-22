package sms

import "context"

type SMS interface {
	SendCode(ctx context.Context, phone, code string) error
}
