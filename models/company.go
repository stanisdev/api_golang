package models

import (
	"github.com/jinzhu/gorm"
	structs "app/structures"
)

type Company struct {
  gorm.Model
	Name string `gorm:"size:150;unique;not null"valid:"required,length(1|150)"`
}

func (dm *DbMethods) CountNotificationsByPublishers(ids []uint) (*[]structs.NotificationsCount, []error) {
	results := &[]structs.NotificationsCount{}
	queryResult := dm.DB.
		Table("notifications").
		Select("COUNT(company_id) AS total, company_id AS publisher_id").
		Where("company_id IN (?)", ids).
		Group("company_id").
		Scan(&results).
		GetErrors()
	return results, queryResult
}