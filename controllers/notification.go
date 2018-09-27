package controllers

import (
	"app/models"
	"app/services"
	"github.com/gin-gonic/gin"
	_ "fmt"
)

func (e *Env) NotificationList(c *gin.Context) {
	ntfs := models.GetDmInstance().FindNotifications() // @TODO: Remove excessive fields
	c.JSON(200, gin.H{
		"ok": true,
		"payload": ntfs,
	})
}

func (e *Env) NotificationCreate(c *gin.Context) {
	cmp := &models.Company{
		Name: c.PostForm("company"),
	}
	if (!models.ValidateModel(cmp)) {
		services.WrongPostData(c)
		return
	}
	exCmp := models.Company{}
	e.db.Where("name = ?", cmp.Name).First(&exCmp)

	var cmpId uint
	if (exCmp.ID < 1) { // No such company, let's create it
		e.db.Create(&cmp)
		cmpId = cmp.ID
	} else {
		cmpId = exCmp.ID
	}

	ntf := &models.Notification{
		Text: c.PostForm("text"),
		Image: c.PostForm("image"),
		CompanyID: cmpId,
	}
	if (!models.ValidateModel(ntf)) {
		services.WrongPostData(c)
		return
	}
	e.db.Create(ntf)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": gin.H{
			"id": ntf.ID,
			"text": ntf.Text,
			"image": ntf.Image,
			"company": cmp.Name,
		},
	})
}

func (e *Env) NotificationRemove(c *gin.Context) {
	ntf := c.MustGet("notification").(*models.NotificationQuery)
	e.db.Where("id = ?", ntf.Id).Limit(1).Unscoped().Delete(&models.Notification{})
	c.JSON(200, gin.H{
		"ok": true,
	})
}

func (e *Env) NotificationGetById(c *gin.Context) {
	ntf := c.MustGet("notification").(*models.NotificationQuery)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": gin.H{
			"id": ntf.Id,
			"text": ntf.Text,
			"image": ntf.Image,
			"company": ntf.Company,
			"created_at": ntf.CreatedAt.Unix(),
		},
	})
}