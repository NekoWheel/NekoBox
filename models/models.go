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
	PageID uint
	Answer string
}
