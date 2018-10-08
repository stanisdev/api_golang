package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"strconv"
)

type Notification struct {
  gorm.Model
	Message string `gorm:"size:3000;not null"valid:"required,length(1|3000)"`
	Image string `gorm:"size:150;not null"valid:"required,length(1|150)"`
	Header string `gorm:"size:200;not null"valid:"required,length(1|200)"`
	Priority uint `gorm:"not null"valid:"required"`
	Expired time.Time
	Button string `gorm:"size:150;not null"valid:"required,length(1|150)"`
	Link string `gorm:"size:250;not null"valid:"required,length(1|250)"`
	CompanyID uint `gorm:"not null"valid:"required"`
}

type NotificationQuery struct {
	Notification
	Company string
}

func (n Notification) GetExpired() string {
	return n.Expired.Format("2006/01/02")
}

func (n Notification) GetExpiredForAll() string {
	day := strconv.Itoa(n.Expired.Day())
	year := strconv.Itoa(n.Expired.Year())
	month := strconv.Itoa(int(n.Expired.Month()))
	return  day + "/" + month + "/" + year
}

func (dm *DbMethods) FindNotificationById(id uint) *NotificationQuery {
	ntf := &NotificationQuery{}
	dm.DB.Table("notifications n").
		Select(`
			n.id,
			n.message,
			n.image,
			n.header,
			n.priority,
			n.expired,
			n.button,
			n.link,
			n.company_id,
			n.created_at,
			c.name company
		`).
		Joins("LEFT JOIN companies c ON n.company_id = c.id").
		Where("n.id = ?", id).
		Order("n.id DESC").
		Limit(1).
		Scan(ntf)
	return ntf
}

func (dm *DbMethods) FindNotifications(message string, limit int, offset int) *[]NotificationQuery {
	ntfs := &[]NotificationQuery{}
	like := "%" + message + "%"
	dm.DB.Table("notifications n").
		Select(`
			n.id,
			n.message,
			n.image,
			n.header,
			n.priority,
			n.expired,
			n.button,
			n.link,
			n.company_id,
			n.created_at,
			c.name company
		`).
		Joins("LEFT JOIN companies c ON n.company_id = c.id").
		Where("n.message LIKE ?", like).
		Order("n.id ASC").
		Limit(limit).
		Offset(offset).
		Scan(ntfs)

	return ntfs
}

func (dm *DbMethods) FindAllNotifications() *[]NotificationQuery {
	ntfs := &[]NotificationQuery{}
	dm.DB.Table("notifications n").
		Select(`
			n.id,
			n.message,
			n.image,
			n.header,
			n.priority,
			n.expired,
			n.button,
			n.link,
			n.company_id,
			n.created_at,
			"" AS name
		`).
		// Joins("LEFT JOIN companies c ON n.company_id = c.id").
		Order("n.id ASC").
		Scan(ntfs)

	return ntfs
}