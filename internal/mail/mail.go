// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mail

import (
	"bytes"
	"crypto/tls"
	"embed"
	"fmt"
	"html/template"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/templates"
)

func SendNewQuestionMail(email string, domain string, questionID uint, questionContent string) error {
	params := map[string]string{
		"link":     fmt.Sprintf("https://box.n3ko.co/_/%s/%d", domain, questionID),
		"question": questionContent,
	}
	return sendTemplateMail(email, "【NekoBox】您有一个新的提问", templates.FS, "mail/new-question.html", params)
}

func SendPasswordRecoveryMail(email, code string) error {
	params := map[string]string{
		"link":  fmt.Sprintf("https://box.n3ko.co/recover-password?code=%s", code),
		"email": email,
	}
	return sendTemplateMail(email, "【NekoBox】账号密码找回", templates.FS, "mail/password-recovery.html", params)
}

func sendTemplateMail(email, title string, templateFS embed.FS, templatePath string, params map[string]string) error {
	var content bytes.Buffer
	t, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		return errors.Wrap(err, "parse template file")
	}
	if err := t.Execute(&content, params); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return sendMail(email, title, content.String())
}

func sendMail(to, title, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Mail.Account)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(
		conf.Mail.SMTP,
		conf.Mail.Port,
		conf.Mail.Account,
		conf.Mail.Password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}
