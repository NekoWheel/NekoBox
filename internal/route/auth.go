package route

import (
	"net/http"

	"github.com/NekoWheel/NekoBox/internal/response"
	"github.com/flamego/recaptcha"
	"github.com/flamego/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/form"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (*AuthHandler) SignUp(ctx context.Context, recaptcha recaptcha.RecaptchaV3, f form.SignUp) error {
	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
		return ctx.Error(http.StatusInternalServerError, "无感验证码请求失败，请稍后再试")
	}
	if !resp.Success {
		return ctx.Error(http.StatusBadRequest, "无感验证码校验失败，请重试")
	}

	if err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		Name:       f.Name,
		Password:   f.Password,
		Email:      f.Email,
		Avatar:     conf.Upload.DefaultAvatarURL,
		Domain:     f.Domain,
		Background: conf.Upload.DefaultBackground,
		Intro:      "问你想问的",
	}); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotExists),
			errors.Is(err, db.ErrBadCredential),
			errors.Is(err, db.ErrDuplicateEmail),
			errors.Is(err, db.ErrDuplicateDomain):
			return ctx.Error(http.StatusBadRequest, errors.Cause(err).Error())

		default:
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create new user")
			return ctx.ServerError()
		}
	}

	return ctx.Success("注册成功，欢迎来到 NekoBox！")
}

func (*AuthHandler) SignIn(ctx context.Context, sess session.Session, recaptcha recaptcha.RecaptchaV3, f form.SignIn) error {
	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
		return ctx.Error(http.StatusInternalServerError, "无感验证码请求失败，请稍后再试")
	}
	if !resp.Success {
		return ctx.Error(http.StatusBadRequest, "无感验证码校验失败，请重试")
	}

	user, err := db.Users.Authenticate(ctx.Request().Context(), f.Email, f.Password)
	if err != nil {
		if errors.Is(err, db.ErrBadCredential) {
			return ctx.Error(http.StatusBadRequest, "电子邮箱或密码错误")
		}
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to authenticate user")
		return ctx.ServerError()
	}

	sess.Set(context.SessionKeyUserID, user.ID)

	return ctx.Success(response.SignIn{
		Profile: &response.SignInUserProfile{
			UID:    user.UID,
			Name:   user.Name,
			Domain: user.Domain,
		},
		SessionID: sess.ID(),
	})
}
