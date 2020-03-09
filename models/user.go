package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

func Register(form *UserRegisterForm) error {
	// check name
	var count int
	DB.Model(&User{}).Where(&User{Email: form.Email}).Count(&count)
	if count != 0 {
		return errors.New("这个邮箱已经注册过账号了！")
	}

	DB.Model(&User{}).Where(&User{Name: form.Name}).Count(&count)
	if count != 0 {
		return errors.New("昵称重复了，换一个吧~")
	}

	DB.Model(&Page{}).Where(&Page{Domain: form.Domain}).Count(&count)
	if count != 0 {
		return errors.New("个性域名重复了，换一个吧~")
	}

	user := new(User)
	user.Name = form.Name
	user.Password = AddSalt(form.Password)
	user.Email = form.Email

	// create page
	page := new(Page)
	page.Domain = form.Domain
	page.Intro = "问你想问的"

	tx := DB.Begin()
	if tx.Create(&page).RowsAffected != 1 {
		tx.Rollback()
		return errors.New("注册失败，好像是服务器坏了...")
	}
	user.PageID = page.ID

	if tx.Create(&user).RowsAffected != 1 {
		tx.Rollback()
		return errors.New("注册失败，好像是服务器坏了...")
	}
	tx.Commit()
	return nil
}

func Login(form *UserLoginForm) (*User, error) {
	user := new(User)
	DB.Model(&User{}).Where(&User{Email: form.Email}).Find(&user)
	if user.Email == "" {
		return &User{}, errors.New("")
	}

	if user.Password == AddSalt(form.Password) {
		return user, nil
	}

	return &User{}, errors.New("")
}

func GetUserByPage(pageId uint) (*User, error) {
	user := new(User)
	DB.Model(&User{}).Where(&User{PageID: pageId}).Find(&user)
	if user.Name == "" {
		return &User{}, errors.New("")
	}
	return user, nil
}

func UpdateUser(id uint, u *User) {
	tx := DB.Begin()
	if tx.Model(&User{}).Where(&User{Model: gorm.Model{ID: id}}).Update(u).RowsAffected != 1 {
		tx.Rollback()
		return
	}
	tx.Commit()
}
