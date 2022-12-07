// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type FlashType string

const (
	Info    = "info"
	Success = "success"
	Warning = "warning"
	Error   = "error"
)

type Flash struct {
	Type     FlashType
	Message  string
	FlashTip string
}

func (c Context) SetInfoFlash(message string, tip ...string) {
	c.Session.SetFlash(Flash{Type: Info, Message: message, FlashTip: strings.Join(tip, "")})
}

func (c Context) SetSuccessFlash(message string, tip ...string) {
	c.Session.SetFlash(Flash{Type: Success, Message: message, FlashTip: strings.Join(tip, "")})
}

func (c Context) SetWarningFlash(message string, tip ...string) {
	c.Session.SetFlash(Flash{Type: Warning, Message: message, FlashTip: strings.Join(tip, "")})
}

func (c Context) SetErrorFlash(message string, tip ...string) {
	c.Session.SetFlash(Flash{Type: Error, Message: message, FlashTip: strings.Join(tip, "")})
}

func (c Context) SetInternalErrorFlash() {
	span := trace.SpanFromContext(c.Request().Context())
	traceID := span.SpanContext().TraceID()

	c.Session.SetFlash(Flash{Type: Error, Message: "服务内部错误，请稍后重试。", FlashTip: fmt.Sprintf("若问题一直出现，请带上该段字符 %s 提交反馈。", traceID.String())})
}
