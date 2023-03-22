// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/flamego/csrf"
	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"github.com/flamego/template"
	"github.com/pkg/errors"
	"github.com/unknwon/com"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/security/sms"
	templatepkg "github.com/NekoWheel/NekoBox/internal/template"
)

type EndpointType string

const (
	EndpointAPI EndpointType = "api"
	EndpointWeb EndpointType = "web"
)

func (e EndpointType) IsAPI() bool {
	return e == EndpointAPI
}

func (e EndpointType) IsWeb() bool {
	return e == EndpointWeb
}

func APIEndpoint(ctx Context) {
	ctx.Map(EndpointAPI)
}

// Context represents context of a request.
type Context struct {
	flamego.Context

	Data     template.Data
	Session  session.Session
	Template template.Template

	User     *db.User
	IsLogged bool
}

// HasError returns true if error occurs in form validation.
func (c *Context) HasError() bool {
	hasErr, ok := c.Data["HasError"]
	if !ok {
		return false
	}
	return hasErr.(bool)
}

func (c *Context) SetError(err error, f ...interface{}) {
	c.Data["HasError"] = true
	c.Data["Error"] = err.Error()

	// Set back the form data.
	if len(f) > 0 {
		form := f[0]
		typ := reflect.TypeOf(form)
		val := reflect.ValueOf(form)

		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			val = val.Elem()
		}

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldName := com.ToSnakeCase(field.Name)

			c.Data[fieldName] = val.Field(i).Interface()
		}
	}

	span := trace.SpanFromContext(c.Request().Context())
	if span.IsRecording() {
		span.SetAttributes(
			attribute.String("nekobox.error", err.Error()),
		)
	}
}

func (c *Context) SetInternalError(f ...interface{}) {
	span := trace.SpanFromContext(c.Request().Context())
	traceID := span.SpanContext().TraceID()

	c.Data["FlashTip"] = fmt.Sprintf("若问题一直出现，请带上该段字符 %s 提交反馈。", traceID.String())
	c.SetError(errors.New("服务内部错误，请稍后重试。"), f...)
}

// Success renders HTML template with given name with 200 OK status code.
func (c *Context) Success(templateName string) {
	c.Template.HTML(http.StatusOK, templateName)
}

func (c *Context) SetTitle(title string) {
	c.Data["Title"] = title
}

func (c *Context) Refresh() {
	c.Redirect(c.Request().URL.Path)
}

func (c *Context) JSON(data interface{}) error {
	resp := map[string]interface{}{
		"code":    0,
		"data":    data,
		"message": "success",
	}

	return json.NewEncoder(c.ResponseWriter()).Encode(resp)
}

func (c *Context) ServerError() error {
	return c.JSONError(50000, "服务器内部错误")
}

func (c *Context) JSONError(errorCode int, message string) error {
	span := trace.SpanFromContext(c.Request().Context())

	resp := map[string]interface{}{
		"code":     errorCode,
		"message":  message,
		"trace_id": span.SpanContext().TraceID().String(),
	}

	statusCode := errorCode / 100
	if statusCode < 100 || statusCode > 999 {
		statusCode = http.StatusInternalServerError
	}
	c.ResponseWriter().WriteHeader(statusCode)

	return json.NewEncoder(c.ResponseWriter()).Encode(resp)
}

// Contexter initializes a classic context for a request.
func Contexter() flamego.Handler {
	return func(ctx flamego.Context, data template.Data, session session.Session, x csrf.CSRF, t template.Template, flash session.Flash) {
		c := Context{
			Context:  ctx,
			Data:     data,
			Session:  session,
			Template: t,
		}

		if ctx.Request().Method == http.MethodPost {
			x.Validate(ctx)
		}

		// Get user from session or header when possible
		c.User = authenticatedUser(c.Context, c.Session)

		var userID uint
		if c.User != nil {
			c.IsLogged = true
			c.Data["IsLogged"] = c.IsLogged
			c.Data["LoggedUser"] = c.User
			c.Data["LoggedUserID"] = c.User.ID
			c.Data["LoggedUserName"] = c.User.Name

			userID = c.User.ID
		} else {
			c.Data["LoggedUserID"] = 0
			c.Data["LoggedUserName"] = ""
		}

		span := trace.SpanFromContext(ctx.Request().Context())
		if span.IsRecording() {
			span.SetAttributes(
				attribute.Bool("nekobox.user.is-login", c.IsLogged),
				attribute.Int("nekobox.user.id", int(userID)),
			)
		}
		c.ResponseWriter().Header().Set("Trace-ID", span.SpanContext().TraceID().String())

		if flash != nil {
			flash, ok := flash.(Flash)
			if ok {
				switch flash.Type {
				case "success":
					c.Data["Success"] = flash.Message
				case "error":
					c.Data["HasError"] = true
					c.Data["Error"] = flash.Message
				case "info":
					c.Data["Info"] = flash.Message
				case "warning":
					c.Data["Warning"] = flash.Message
				}
			}
		}

		c.SetTitle("NekoBox")
		c.Data["CSRFToken"] = x.Token()
		c.Data["CSRFTokenHTML"] = templatepkg.Safe(`<input type="hidden" name="_csrf" value="` + x.Token() + `">`)

		c.Data["RecaptchaDomain"] = conf.Recaptcha.Domain
		c.Data["RecaptchaSiteKey"] = conf.Recaptcha.SiteKey
		c.Data["CurrentURI"] = ctx.Request().Request.RequestURI
		c.Data["ExternalURL"] = conf.App.ExternalURL

		// 🚨 SECURITY: Prevent MIME type sniffing in some browsers,
		c.ResponseWriter().Header().Set("X-Content-Type-Options", "nosniff")
		c.ResponseWriter().Header().Set("X-Frame-Options", "DENY")

		var smsModule sms.SMS
		if conf.SMS.AliyunSignName != "" && conf.SMS.AliyunTemplateCode != "" {
			smsModule = sms.NewAliyunSMS(sms.NewAliyunSMSOptions{
				Region:          conf.SMS.AliyunRegion,
				AccessKey:       conf.SMS.AliyunAccessKey,
				AccessKeySecret: conf.SMS.AliyunAccessKeySecret,
				SignName:        conf.SMS.AliyunSignName,
				TemplateCode:    conf.SMS.AliyunTemplateCode,
			})
		} else {
			smsModule = sms.NewDummySMS()
		}
		ctx.MapTo(smsModule, (sms.SMS)(nil))

		ctx.Map(c)
		ctx.Map(EndpointWeb)
	}
}
