// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"encoding/json"
	"html"
	"html/template"
	"strings"
	"sync"
	"time"

	"gorm.io/datatypes"

	"github.com/wuhan005/NekoBox/internal/conf"
)

var (
	funcMap     []template.FuncMap
	funcMapOnce sync.Once
)

func FuncMap() []template.FuncMap {
	funcMapOnce.Do(func() {
		funcMap = []template.FuncMap{map[string]interface{}{
			"ICP": func() string {
				return conf.App.ICP
			},
			"CommitSHA": func() string {
				return conf.BuildCommit
			},
			"CommitSHAShort": func() string {
				if len(conf.BuildCommit) > 7 {
					return conf.BuildCommit[:7]
				}
				return conf.BuildCommit
			},
			"Date": func(t time.Time, format string) string {
				replacer := strings.NewReplacer(datePatterns...)
				format = replacer.Replace(format)
				return t.Format(format)
			},
			"QuestionFormat": func(input string) template.HTML {
				input = html.EscapeString(input)
				input = strings.ReplaceAll(input, "\n", "</br>")
				return template.HTML(input)
			},
			"AnswerFormat": func(input string) template.HTML {
				input = html.EscapeString(input)
				input = strings.ReplaceAll(input, "\n", "</br>")
				return template.HTML(input)
			},
			"SentryDSN": func() string {
				return conf.App.SentryDSN
			},
			"ParsePublicURLs": func(input datatypes.JSON) string {
				urls := make(map[string]string)
				if err := json.Unmarshal(input, &urls); err != nil {
					return ""
				}
				for _, url := range urls {
					return url
				}
				return ""
			},
			"ImageBucketCDNHost": func() string {
				return conf.Upload.ImageBucketCDNHost
			},
			"Safe": Safe,
		}}
	})
	return funcMap
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}
