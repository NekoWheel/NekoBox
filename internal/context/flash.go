// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"

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
	Type    FlashType
	Message string
}

func (c Context) SetInfoFlash(message string) {
	c.Session.SetFlash(Flash{Type: Info, Message: message})
}

func (c Context) SetSuccessFlash(message string) {
	c.Session.SetFlash(Flash{Type: Success, Message: message})
}

func (c Context) SetWarningFlash(message string) {
	c.Session.SetFlash(Flash{Type: Warning, Message: message})
}

func (c Context) SetErrorFlash(message string) {
	c.Session.SetFlash(Flash{Type: Error, Message: message})
}

func (c Context) SetInternalErrorFlash() {
	span := trace.SpanFromContext(c.Request().Context())
	traceID := span.SpanContext().TraceID()

	c.Session.SetFlash(Flash{Type: Error, Message: fmt.Sprintf("服务内部错误，请稍后重试。若问题一直出现，请带上该段字符 %s 提交反馈。", traceID.String())})
}
