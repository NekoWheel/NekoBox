// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"net/http"

	"github.com/flamego/csrf"
	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"github.com/flamego/template"
	"gorm.io/gorm"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/dbutil"
	templatepkg "github.com/NekoWheel/NekoBox/internal/template"
)

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

func (c *Context) SetError(err error) {
	c.Data["HasError"] = true
	c.Data["Error"] = err.Error()
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

// Contexter initializes a classic context for a request.
func Contexter(db *gorm.DB) flamego.Handler {
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

		if c.User != nil {
			c.IsLogged = true
			c.Data["IsLogged"] = c.IsLogged
			c.Data["LoggedUser"] = c.User
			c.Data["LoggedUserID"] = c.User.ID
			c.Data["LoggedUserName"] = c.User.Name
		} else {
			c.Data["LoggedUserID"] = 0
			c.Data["LoggedUserName"] = ""
		}

		// If request sends files, parse them here otherwise the Query() can't be parsed and the CsrfToken will be invalid.
		//if c.Request().Method == http.MethodPost && strings.Contains(c.Request().Header.Get("Content-Type"), "multipart/form-data") {
		//	if err := c.Request().ParseMultipartForm(conf.Attachment.MaxSize << 20); err != nil && !strings.Contains(err.Error(), "EOF") { // 32MB max size
		//		c.Error(err, "parse multipart form")
		//		return
		//	}
		//}

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

		// 🚨 SECURITY: Prevent MIME type sniffing in some browsers,
		c.ResponseWriter().Header().Set("X-Content-Type-Options", "nosniff")
		c.ResponseWriter().Header().Set("X-Frame-Options", "DENY")

		c.MapTo(db, (*dbutil.Transactor)(nil))
		ctx.Map(c)
	}
}
