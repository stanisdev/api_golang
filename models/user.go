package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
  gorm.Model
	Username string `gorm:"size:50;unique;not null"valid:"required,length(1|50)"`
	Password string `gorm:"size:60;not null"`
	Salt string `gorm:"size:10;not null"`
	UniqUserKey string `gorm:"size:10;unique;not null"valid:"required,length(10)"`
	LastLogin time.Time
	PasswordChanged time.Time
}

func (dm DbMethods) UserFindOne(params *User) (*User, bool) {
	user := User{}
	instance.Where(params).First(&user)
	return &user, user.ID > 0
}