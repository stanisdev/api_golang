package controllers

import (
	"app/models"
	structs "app/structures"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "fmt"
)

func (e *Env) NotificationList(c *gin.Context) {
	text := c.Query("text")
	lmt := c.MustGet("limit").(int)
	ofst := c.MustGet("offset").(int)
	
	ntfs := models.GetDmInstance().FindNotifications(text, lmt, ofst)
	var result []structs.NotificationList
	var msg string
	for _, ntf := range *ntfs {
		msg = ntf.Message
		if (len(msg) > 62) {
			msg = msg[0:62] + "..."
		}
		result = append(result, structs.NotificationList {
			Id: ntf.ID,
			Message: msg,
			Expired: ntf.Expired.Format("Jan 2 2006"),
			Link: ntf.Link,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"payload": result,
	})
}

func (e *Env) NotificationCreate(c *gin.Context) {
	ntf := c.MustGet("notificationBlank").(*models.Notification)
	e.db.Create(ntf)
	c.JSON(200, gin.H{
		"ok": true,
	})
}

func (e *Env) NotificationUpdate(c *gin.Context) {
	ntfQuery := c.MustGet("notification").(*models.NotificationQuery)
	ntfBlank := c.MustGet("notificationBlank").(*models.Notification)
	ntfBlank.ID = ntfQuery.ID
	ntfBlank.CreatedAt = ntfQuery.CreatedAt 

	e.db.Save(&ntfBlank)

	oldCompanyID := c.MustGet("oldCompanyID").(uint)
	var count int
	e.db.Model(&models.Notification{}).Where("company_id = ?", oldCompanyID).Count(&count)
	if (count < 1) {
		e.db.Where("id = ?", oldCompanyID).Limit(1).Unscoped().Delete(&models.Company{}) // Remove outdated single company
	}
	c.JSON(200, gin.H{
		"ok": true,
	})
}

func (e *Env) NotificationRemove(c *gin.Context) {
	ntf := c.MustGet("notification").(*models.NotificationQuery)
	e.db.Where("id = ?", ntf.ID).Limit(1).Unscoped().Delete(&models.Notification{})

	var count int
	e.db.Model(&models.Notification{}).Where("company_id = ?", ntf.CompanyID).Count(&count)
	if (count < 1) {
		e.db.Where("id = ?", ntf.CompanyID).Limit(1).Unscoped().Delete(&models.Company{})
	}
	c.JSON(200, gin.H{
		"ok": true,
	})
}

func (e *Env) NotificationGetById(c *gin.Context) {
	ntf := c.MustGet("notification").(*models.NotificationQuery)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": gin.H{
			"id": ntf.ID,
			"message": ntf.Message,
			"image": ntf.Image,
			"header": ntf.Header,
			"company": ntf.Company,
			"priority": ntf.Priority,
			"expired": ntf.Expired.Unix(),
			"button": ntf.Button,
			"link": ntf.Link,
			"created_at": ntf.CreatedAt.Unix(),
		},
	})
}

func (e *Env) NotificationCount(c *gin.Context) {
	var count int
	e.db.Model(&models.Notification{}).Count(&count)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": count,
	})
}

func (e *Env) NotificationPublic(c *gin.Context) {
	var result []structs.NotificationPublic
	ntfs := models.GetDmInstance().FindAllNotifications()
	for _, ntf := range *ntfs {
		result = append(result, structs.NotificationPublic {
			Id: ntf.ID,
			Image: ntf.Image,
			Message: ntf.Message,
			Header: ntf.Header,
			Priority: ntf.Priority,
			Expired: ntf.GetExpiredForAll(),
			Button: ntf.Button,
			Link: ntf.Link,
		})
	}
	c.JSON(200, gin.H{
		"ok": true,
		"payload": result,
	})
}