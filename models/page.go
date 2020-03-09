package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

func GetPage(id uint) (*Page, error) {
	page := new(Page)
	DB.Model(&Page{}).Where(&Page{Model: gorm.Model{ID: id}}).Find(&page)

	if page.ID == 0 {
		return &Page{}, errors.New("服务器错误！")
	}
	return page, nil
}
