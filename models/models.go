package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

func init() {
	validation.SetDefaultMessage(map[string]string{
		"Required": "不能为空",
		"MinSize":  "长度最小值是 %d",
		"MaxSize":  "长度最大值是 %d",
		"Length":   "长度需要为 %d",
		"Email":    "电子邮箱格式不正确",
	})

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True",
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
	Name     string `form:"name" valid:"Required; MaxSize(20)"`
	Password string `form:"password" valid:"Required; MinSize(8); MaxSize(30)"`
	Email    string `form:"email" valid:"Required; Email; MaxSize(100)"`
	Domain   string `form:"domain" valid:"Required; AlphaDash; MinSize(3); MaxSize(10)"`
}

type UserLoginForm struct {
	Email    string `form:"email" valid:"Required; Email; MaxSize(100)"`
	Password string `form:"password" valid:"Required; MinSize(8); MaxSize(30)"`
}

type QuestionForm struct {
	PageID  uint
	Content string `form:"content" valid:"Required; MaxSize(50)"`
}

type UpdateForm struct {
	Name     string `form:"name" valid:"Required; MaxSize(20)"`
	Password string `form:"password"`
	Intro    string `form:"intro" valid:"MaxSize(40)"`
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
