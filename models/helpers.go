package models

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"github.com/astaxie/beego"
	"io"
)

func addSalt(raw string) string {
	return hmacSha1Encode(raw, beego.AppConfig.String("salt"))
}

func hmacSha1Encode(input string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	_, _ = io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}
