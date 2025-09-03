package route

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/flamego/recaptcha"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/wuhan005/govalid"
	"gorm.io/gorm"

	"github.com/NekoWheel/NekoBox/internal/conf"
	"github.com/NekoWheel/NekoBox/internal/context"
	"github.com/NekoWheel/NekoBox/internal/db"
	"github.com/NekoWheel/NekoBox/internal/dbutil"
	"github.com/NekoWheel/NekoBox/internal/form"
	"github.com/NekoWheel/NekoBox/internal/mail"
	"github.com/NekoWheel/NekoBox/internal/response"
	"github.com/NekoWheel/NekoBox/internal/security/censor"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (*UserHandler) Domainer(ctx context.Context) error {
	domain := ctx.Param("domain")
	pageUser, err := db.Users.GetByDomain(ctx.Request().Context(), domain)
	if err != nil {
		if errors.Is(err, db.ErrUserNotExists) {
			return ctx.Error(http.StatusNotFound, "提问箱不存在")
		}

		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get user by domain")
		return ctx.ServerError()
	}

	if pageUser.HarassmentSetting == db.HarassmentSettingTypeRegisterOnly && !ctx.IsSignedIn {
		return ctx.Error(http.StatusUnauthorized, "请先登录账号")
	}

	ctx.Map(pageUser)
	return nil
}

func (*UserHandler) OwnerRequired(ctx context.Context, pageUser *db.User) error {
	if ctx.IsSignedIn && pageUser.ID == ctx.User.ID {
		return nil
	}
	return ctx.Error(http.StatusForbidden, "无权访问")
}

func (*UserHandler) Profile(ctx context.Context, pageUser *db.User) error {
	return ctx.Success(&response.UserProfile{
		UID:               pageUser.UID,
		Name:              pageUser.Name,
		Avatar:            pageUser.Avatar,
		Domain:            pageUser.Domain,
		Background:        pageUser.Background,
		Intro:             pageUser.Intro,
		HarassmentSetting: string(pageUser.HarassmentSetting),
	})
}

func (*UserHandler) ListQuestions(ctx context.Context, pageUser *db.User) error {
	pageSize := ctx.QueryInt("pageSize")
	cursorValue := ctx.Query("cursor")

	total, err := db.Questions.Count(ctx.Request().Context(), pageUser.ID, db.GetQuestionsCountOptions{
		FilterAnswered: true,
		ShowPrivate:    false,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get questions count")
		return ctx.ServerError()
	}

	questions, err := db.Questions.GetByUserID(ctx.Request().Context(), pageUser.ID, db.GetQuestionsByUserIDOptions{
		Cursor: &dbutil.Cursor{
			Value:    cursorValue,
			PageSize: pageSize,
		},
		FilterAnswered: true,
		ShowPrivate:    false,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get questions by user ID")
		return ctx.ServerError()
	}

	respQuestions := lo.Map(questions, func(question *db.Question, _ int) *response.PageQuestionsItem {
		return &response.PageQuestionsItem{
			ID:        question.ID,
			CreatedAt: question.CreatedAt,
			Content:   question.Content,
			Answer:    question.Answer,
		}
	})

	var cursor string
	if len(questions) > 0 {
		cursor = strconv.Itoa(int(questions[len(questions)-1].ID))
	}

	return ctx.Success(&response.PageQuestions{
		Total:     total,
		Cursor:    cursor,
		Questions: respQuestions,
	})
}

func (*UserHandler) PostQuestion(ctx context.Context, pageUser *db.User, recaptcha recaptcha.RecaptchaV3, tx dbutil.Transactor, f form.PostQuestion) error {
	if !ctx.IsSignedIn && pageUser.HarassmentSetting == db.HarassmentSettingTypeRegisterOnly {
		return ctx.Error(http.StatusBadRequest, "提问箱的主人设置了仅注册用户才能提问，请先登录。")
	}

	var receiveReplyEmail string
	if f.ReceiveReplyEmail != "" {
		// Check the email address is valid.
		if errs, ok := govalid.Check(struct {
			Email string `valid:"required;email" label:"邮箱地址"`
		}{
			Email: f.ReceiveReplyEmail,
		}); !ok {
			//nolint:govet
			return ctx.Error(http.StatusBadRequest, errs[0].Error())
		}
	}

	// Check recaptcha code.
	resp, err := recaptcha.Verify(f.Recaptcha, ctx.Request().Request.RemoteAddr)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to check recaptcha")
		return ctx.Error(http.StatusInternalServerError, "无感验证码请求失败，请稍后再试")
	}
	if !resp.Success {
		return ctx.Error(http.StatusBadRequest, "无感验证码校验失败，请重试")
	}

	content := f.Content
	// 🚨 User's block words check.
	if len(pageUser.BlockWords) > 0 {
		blockWords := strings.Split(pageUser.BlockWords, ",")
		for _, word := range blockWords {
			if strings.Contains(content, word) {
				return ctx.Error(http.StatusBadRequest, "提问内容中包含了提问箱主人设置的屏蔽词，发送失败")
			}
		}
	}

	// 🚨 Content security check.
	censorResponse, err := censor.Text(ctx.Request().Context(), content)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to censor text")
	}
	if err == nil && !censorResponse.Pass {
		errorMessage := censorResponse.ErrorMessage()
		//nolint:govet
		return ctx.Error(http.StatusBadRequest, errorMessage)
	}

	// ⚠️ Here is the aliyun CDN origin IP header.
	// A security problem may occur if the CDN is enabled and users can modify the header.
	fromIP := ctx.Request().Header.Get("Ali-CDN-Real-IP")
	if fromIP == "" {
		fromIP = ctx.Request().Header.Get("CF-Connecting-IP")
	}
	if fromIP == "" {
		fromIP = ctx.Request().Header.Get("X-Real-IP")
	}

	// Try to get current logged user.
	var askerUserID uint
	if ctx.IsSignedIn {
		askerUserID = ctx.User.ID
	}

	// Upload image if exists.
	var uploadImage *db.UploadImage
	if len(f.Images) > 0 {
		image := f.Images[0]
		uploadImage, err = uploadImageFile(ctx, uploadImageFileOptions{
			Image:          image,
			UploaderUserID: askerUserID,
		})
		if err != nil {
			if errors.Is(err, ErrUploadImageSizeTooLarge) {
				return ctx.Error(http.StatusBadRequest, "图片文件大小不能大于 5Mb")
			} else {
				logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to upload image")
				return ctx.Error(http.StatusInternalServerError, "上传图片失败，请重试")
			}
		}
	}

	var question *db.Question
	if err := tx.Transaction(func(tx *gorm.DB) error {
		questionsStore := db.NewQuestionsStore(tx)

		question, err = questionsStore.Create(ctx.Request().Context(), db.CreateQuestionOptions{
			FromIP:            fromIP,
			UserID:            pageUser.ID,
			Content:           content,
			ReceiveReplyEmail: receiveReplyEmail,
			AskerUserID:       askerUserID,
			IsPrivate:         f.IsPrivate,
		})
		if err != nil {
			return errors.Wrap(err, "create question")
		}

		// Update censor result.
		if err := questionsStore.UpdateCensor(ctx.Request().Context(), question.ID, db.UpdateQuestionCensorOptions{
			ContentCensorMetadata: censorResponse.ToJSON(),
		}); err != nil {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to update question censor result")
			return errors.Wrap(err, "update question censor")
		}

		if uploadImage != nil {
			// Bind the uploaded image with the question.
			if err := db.NewUploadImagesStore(tx).BindUploadImageWithQuestion(ctx.Request().Context(), uploadImage.ID, db.UploadImageQuestionTypeAsk, question.ID); err != nil {
				return errors.Wrap(err, "bind upload image with question")
			}
		}
		return nil
	}); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to post question")
		return ctx.ServerError()
	}

	go func() {
		if pageUser.Notify == db.NotifyTypeEmail {
			// Send notification to page user.
			if err := mail.SendNewQuestionMail(pageUser.Email, pageUser.Domain, question.ID, question.Content); err != nil {
				logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to send new question mail to user")
			}
		}
	}()

	questionPrivateURL := fmt.Sprintf("/_/%s/%d?t=%s", pageUser.Domain, question.ID, question.Token)
	questionPrivateAbsURL := fmt.Sprintf("%s%s", strings.TrimRight(conf.App.ExternalURL, "/"), questionPrivateURL)
	return ctx.Success("发送问题成功！以下是提问私密链接，使用该链接可以随时查看你的提问，请注意保存。" + fmt.Sprintf(`<a href="%s" target="_blank">%[1]s</a>`, questionPrivateAbsURL))
}

type uploadImageFileOptions struct {
	Image          *multipart.FileHeader
	UploaderUserID uint
}

var ErrUploadImageSizeTooLarge = errors.New("图片文件大小不能大于 5Mb")

func uploadImageFile(ctx context.Context, options uploadImageFileOptions) (*db.UploadImage, error) {
	image := options.Image
	fileName := image.Filename
	fileExt := filepath.Ext(fileName)
	fileSize := image.Size
	if fileSize > 1024*1024*5 { // 5Mib
		return nil, ErrUploadImageSizeTooLarge
	}

	now := time.Now()
	fileKey := fmt.Sprintf("%d/%d/%d%s", now.Year(), now.Month(), now.UnixNano(), fileExt)

	uploadImageFile, err := image.Open()
	if err != nil {
		return nil, errors.Wrap(err, "open image")
	}
	defer func() { _ = uploadImageFile.Close() }()

	hasher := md5.New()
	reader := io.TeeReader(uploadImageFile, hasher)

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               conf.Upload.ImageEndpoint,
			HostnameImmutable: true,
			Source:            aws.EndpointSourceCustom,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx.Request().Context(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.Upload.ImageAccessID, conf.Upload.ImageAccessSecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "load config")
	}

	client := s3.NewFromConfig(cfg)
	if _, err := client.PutObject(ctx.Request().Context(), &s3.PutObjectInput{
		Bucket:        aws.String(conf.Upload.ImageBucket),
		Key:           aws.String(fileKey),
		Body:          reader,
		ContentLength: aws.Int64(fileSize),
	}); err != nil {
		return nil, errors.Wrap(err, "put object")
	}

	fileMd5 := fmt.Sprintf("%x", hasher.Sum(nil))

	uploadImage, err := db.UploadImgaes.Create(ctx.Request().Context(), db.CreateUploadImageOptions{
		UploaderUserID: options.UploaderUserID,
		Name:           fileName,
		FileSize:       fileSize,
		Md5:            fileMd5,
		Key:            fileKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create upload image")
	}
	return uploadImage, nil
}

func (*UserHandler) GetQuestion(ctx context.Context, pageUser *db.User) error {
	questionID := uint(ctx.ParamInt("questionID"))
	questionToken := ctx.Query("t")

	question, err := db.Questions.GetByID(ctx.Request().Context(), questionID)
	if err != nil {
		if errors.Is(err, db.ErrQuestionNotExist) {
			return ctx.Error(http.StatusNotFound, "提问不存在")
		}
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get question by ID")
		return ctx.ServerError()
	}

	// Check the question is belongs to the correct page user.
	if question.UserID != pageUser.ID {
		return ctx.Error(http.StatusNotFound, "提问不存在")
	}

	// If the question has not been answered, we should check the question is belongs to the correct page user.
	// The questioner can use the token to view the question.
	if (question.Answer == "" || question.IsPrivate) && (!ctx.IsSignedIn || ctx.User.ID != question.UserID) && (question.Token != "" && question.Token != questionToken) {
		return ctx.Error(http.StatusNotFound, "提问不存在")
	}

	questionUploadImages, err := db.UploadImgaes.GetByTypeQuestionID(ctx.Request().Context(), db.UploadImageQuestionTypeAsk, questionID)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get ask upload images")
		return ctx.ServerError()
	}
	questionImageURLs := lo.Map(questionUploadImages, func(item *db.UploadImage, _ int) string {
		return "https://" + conf.Upload.ImageBucketCDNHost + "/" + item.Key
	})

	answerUploadImages, err := db.UploadImgaes.GetByTypeQuestionID(ctx.Request().Context(), db.UploadImageQuestionTypeAnswer, questionID)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get answer upload images")
		return ctx.ServerError()
	}
	answerImageURLs := lo.Map(answerUploadImages, func(item *db.UploadImage, _ int) string {
		return "https://" + conf.Upload.ImageBucketCDNHost + "/" + item.Key
	})

	return ctx.Success(&response.PageQuestion{
		ID:                question.ID,
		IsOwner:           ctx.IsSignedIn && ctx.User.ID == question.UserID,
		CreatedAt:         question.CreatedAt,
		AnsweredAt:        question.UpdatedAt,
		Content:           question.Content,
		Answer:            question.Answer,
		QuestionImageURLs: questionImageURLs,
		AnswerImageURLs:   answerImageURLs,
		HasReplyEmail:     question.ReceiveReplyEmail != "",
		IsPrivate:         question.IsPrivate,
	})
}
