package models

import (
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	validation.SetDefaultMessage(map[string]string{
		"Required":  "不能为空",
		"MinSize":   "长度最小值是 %d",
		"MaxSize":   "长度最大值是 %d",
		"Length":    "长度需要为 %d",
		"Email":     "格式不正确",
		"AlphaDash": "只能包含字符或数字或横杠 -_",
	})

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local",
			beego.AppConfig.String("db_user"),
			beego.AppConfig.String("db_password"),
			beego.AppConfig.String("db_addr"),
			beego.AppConfig.String("db_name"),
		))

	if err != nil {
		log.Fatalln(err)
	}
	DB = db

	DB.AutoMigrate(&User{}, &Page{}, &Question{})
}

type UserRegisterForm struct {
	Recaptcha      string `form:"g-recaptcha-response" valid:"Required" label:"Recaptcha"`
	Name           string `form:"name" valid:"Required; MaxSize(20)" label:"昵称"`
	Password       string `form:"password" valid:"Required; MinSize(8); MaxSize(30)" label:"密码"`
	RepeatPassword string `form:"repeat_password"`
	Email          string `form:"email" valid:"Required; Email; MaxSize(100)" label:"电子邮箱"`
	Domain         string `form:"domain" valid:"Required; AlphaDash; MinSize(3); MaxSize(10)" label:"个性域名"`
}

func (f *UserRegisterForm) Valid(v *validation.Validation) {
	if f.Password != f.RepeatPassword {
		_ = v.SetError("Password", "两次输入的密码不相同")
	}
}

type UserLoginForm struct {
	Recaptcha string `form:"g-recaptcha-response" valid:"Required" label:"Recaptcha"`
	Email     string `form:"email" valid:"Required; Email; MaxSize(100)" label:"电子邮箱"`
	Password  string `form:"password" valid:"Required; MinSize(8); MaxSize(30)" label:"密码"`
}

type QuestionForm struct {
	Recaptcha string `form:"g-recaptcha-response" valid:"Required" label:"Recaptcha"`
	PageID    uint
	Content   string `form:"content" valid:"Required; MaxSize(50)" label:"问题内容"`
}

type UpdateForm struct {
	Name     string `form:"name" valid:"Required; MaxSize(20)" label:"昵称"`
	Password string `form:"password" label:"密码"`
	Intro    string `form:"intro" valid:"MaxSize(40)" label:"留言板介绍"`
}

type AnswerForm struct {
	Answer string `form:"answer" valid:"Required; MaxSize(150)" label:"回答内容"`
}

type UploadCallBack struct {
	Code int `json:"code"`
	Data struct {
		Md5      string `json:"md5"`
		Mime     string `json:"mime"`
		Name     string `json:"name"`
		Quota    string `json:"quota"`
		Sha1     string `json:"sha1"`
		Size     int    `json:"size"`
		URL      string `json:"url"`
		UseQuota string `json:"use_quota"`
	} `json:"data"`
	Msg  string `json:"msg"`
	Time int    `json:"time"`
}

type RecaptchaResponse struct {
	Success bool `json:"success"`
}

type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string
	Avatar   string
	PageID   uint
}

type Page struct {
	gorm.Model
	Domain     string
	Background string
	Intro      string
}

type Question struct {
	gorm.Model
	PageID  uint
	Content string
	Answer  string
}
