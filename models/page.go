package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

func GetPageByID(id uint) (*Page, error) {
	page := new(Page)
	DB.Model(&Page{}).Where(&Page{Model: gorm.Model{ID: id}}).Find(&page)

	if page.ID == 0 {
		return &Page{}, errors.New("服务器错误！")
	}
	return page, nil
}

func GetPageByDomain(domain string) (*Page, error) {
	page := new(Page)
	DB.Model(&Page{}).Where(&Page{Domain: domain}).Find(page)
	if page.Domain == "" {
		return &Page{}, errors.New("问答箱不存在！")
	}
	return page, nil
}

func UpdatePage(pageID uint, page *Page) {
	tx := DB.Begin()
	if tx.Model(&Page{}).Where(&Page{Model: gorm.Model{ID: pageID}}).Update(page).RowsAffected != 1 {
		tx.Rollback()
		return
	}
	tx.Commit()
}
