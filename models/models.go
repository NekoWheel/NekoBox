package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

func init() {
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

	DB.AutoMigrate()
}
