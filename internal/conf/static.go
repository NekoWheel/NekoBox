// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

// Build time and commit information.
//
// ⚠️ WARNING: should only be set by "-ldflags".
var (
	BuildTime   string
	BuildCommit = "dev"
)

var (
	App struct {
		Production bool   `ini:"production"`
		ICP        string `ini:"icp"`
		UptraceDSN string `ini:"uptrace_dsn"`
	}

	Server struct {
		Port    int    `ini:"port"`
		Salt    string `ini:"salt"`
		XSRFKey string `ini:"xsrf_key"`
	}

	Database struct {
		DSN      string
		User     string `ini:"user"`
		Password string `ini:"password"`
		Address  string `ini:"address"`
		Name     string `ini:"name"`
	}

	Recaptcha struct {
		Domain    string `ini:"domain"`
		SiteKey   string `ini:"site_key"`
		ServerKey string `ini:"server_key"`
	}

	Upload struct {
		Token              string `ini:"token"`
		URL                string `ini:"url"`
		DefaultAvatarURL   string `ini:"default_avatar"`
		DefaultBackground  string `ini:"default_background"`
		AliyunEndpoint     string `ini:"aliyun_endpoint"`
		AliyunAccessID     string `ini:"aliyun_access_id"`
		AliyunAccessSecret string `ini:"aliyun_access_secret"`
		AliyunBucket       string `ini:"aliyun_bucket"`
	}

	Mail struct {
		Account  string `ini:"account"`
		Password string `ini:"password"`
		Port     int    `ini:"port"`
		SMTP     string `ini:"smtp"`
	}
)
