package models

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/parnurzeal/gorequest"
	"io"
	"mime/multipart"
	"strings"
)

func AddSalt(raw string) string {
	return hmacSha1Encode(raw, beego.AppConfig.String("salt"))
}

func hmacSha1Encode(input string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	_, _ = io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CheckRecaptcha(response string, remoteip string) bool {
	req := gorequest.New().Post(beego.AppConfig.String("recaptcha_domain") + "/recaptcha/api/siteverify").Type("form")
	req.SendMap(map[string]string{
		"secret":   beego.AppConfig.String("recaptcha_server_key"),
		"response": response,
		"remoteip": remoteip,
	})
	resp, body, _ := req.End()
	if body == "" || resp == nil || resp.StatusCode != 200 {
		return false
	}

	recaptcha := new(RecaptchaResponse)
	err := json.Unmarshal([]byte(body), &recaptcha)
	if err != nil {
		return false
	}
	if recaptcha.Success {
		return true
	}
	return false
}

func UploadPicture(header *multipart.FileHeader, file multipart.File) string {
	fileByte := make([]byte, header.Size)
	_, _ = file.Read(fileByte)
	req := gorequest.New().Post(beego.AppConfig.String("upload_url")).Type("multipart")
	req.Header.Set("token", beego.AppConfig.String("upload_token"))
	req.SendFile(fileByte, header.Filename, "image")
	resp, body, _ := req.End()

	if resp != nil && resp.StatusCode == 200 {
		backgroundJSON := new(UploadCallBack)
		err := json.Unmarshal([]byte(body), &backgroundJSON)
		if err == nil {
			return strings.Split(backgroundJSON.Data.URL, "?")[0]
		}
	}
	return ""
}
