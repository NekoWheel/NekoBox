package route

import (
	"fmt"
	"os"
	"time"

	"github.com/flamego/cache"
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/security/sms"
	"github.com/NekoWheel/NekoBox/internal/strutil"
)

const (
	SMSCacheKeyRateLimit       = "sms-rate-limit:"
	SMSCacheKeyPrefixRegister  = "register-sms-code:"
	SMSCacheKeyPrefixBindPhone = "bind-phone-sms-code:"
)

func SendSMS(keyPrefix string) func(ctx context.Context, f form.SendSMS, sms sms.SMS, cache cache.Cache, recaptcha recaptcha.RecaptchaV3) error {
	return func(ctx context.Context, f form.SendSMS, sms sms.SMS, cache cache.Cache, recaptcha recaptcha.RecaptchaV3) error {
		// Check recaptcha code.
		resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
		if err != nil {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
			return ctx.ServerError()
		}
		if !resp.Success {
			return ctx.JSONError(40000, "验证码错误")
		}

		if ctx.HasError() {
			return ctx.JSONError(40000, fmt.Sprintf("%s", ctx.Data["Error"]))
		}

		phone := f.Phone

		smsRateLimitCacheKey := SMSCacheKeyRateLimit + phone
		_, err = cache.Get(ctx.Request().Context(), smsRateLimitCacheKey)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to read sms rate limit cache")
			return ctx.ServerError()
		} else if err == nil {
			return ctx.JSONError(40000, "请勿频繁发送短信验证码")
		}

		code := strutil.RandomNumericString(6)

		smsCodeCacheKey := keyPrefix + phone
		if err := cache.Set(ctx.Request().Context(), smsCodeCacheKey, code, 5*time.Minute); err != nil {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to set sms code cache")
			return ctx.ServerError()
		}

		// Set sms rate limit cache.
		if err := cache.Set(ctx.Request().Context(), smsRateLimitCacheKey, time.Now(), 1*time.Minute); err != nil {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to set sms rate limit cache")
			return ctx.ServerError()
		}

		logrus.WithContext(ctx.Request().Context()).
			WithField("key_prefix", keyPrefix).
			WithField("phone", phone).
			WithField("code", code).
			Info("Send sms code successfully")

		if err := sms.SendCode(ctx.Request().Context(), phone, code); err != nil {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to send sms code")
			return ctx.JSONError(50000, "发送短信验证码失败，请稍后重试")
		}

		return ctx.JSON("发送短信验证码成功")
	}
}
