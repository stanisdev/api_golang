package controllers

import (
	"app/models"
	"app/services"
	"github.com/gin-gonic/gin"
	structs "app/structures"
	"net/http"
	_ "fmt"
)

/**
 * Create publisher
 */
func (e *Env) PublisherCreate(c *gin.Context) {
	var publ structs.Publisher
	c.BindJSON(&publ)
	if (!e.ValidateStruct(publ)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
		})
		return
	}
	company := models.Company{}
	result := e.db.Where("name = ?", publ.Name).First(&company).GetErrors()
	if (models.HasError(result)) {
		e.ServerError(c)
		return
	}
	if (company.ID > 0) { // Publisher exists
		c.JSON(http.StatusConflict, gin.H{
			"ok": false,
		})
		return
	}
	newCompany := &models.Company{ 
		Name: publ.Name, 
	}
	result = e.db.Create(newCompany).GetErrors()
	if (models.HasError(result)) {
		e.ServerError(c)
		return
	}
	e.Json(c.Writer, map[string]interface{} {
		"ok": true,
		"payload": gin.H{
			"id": newCompany.ID,
			"name": newCompany.Name,
		},
	})
}

/**
 * Edit publisher
 */
func (e *Env) PublisherUpdate(c *gin.Context) {
	var publ structs.Publisher
	c.BindJSON(&publ)
	if (!e.ValidateStruct(publ)) {
		services.WrongPostData(c)
		return
	}
	publisher := c.MustGet("publisher").(*models.Company)
	if (publisher.Name == publ.Name) {
		e.Json(c.Writer, map[string]interface{} {
			"ok": true,
		})
		return
	}
	_publ := &models.Company{} // Check whether such a publisher name already exists
	result := e.db.Where("name = ?", publ.Name).First(&_publ).GetErrors()
	if (models.HasError(result)) {
		e.ServerError(c)
		return
	}
	if (_publ.ID > 0) {
		c.JSON(200, gin.H{
			"ok": false,
			"errors": gin.H{
				"name": "A company with this name already exists",
			},
		})
		return
	}
	publisher.Name = publ.Name // Update
	result = e.db.Save(&publisher).GetErrors()
	if (models.HasError(result)) {
		e.ServerError(c)
		return
	}
	e.Json(c.Writer, map[string]interface{} {
		"ok": true,
	})
}

/**
 * Remove publisher
 */
func (e *Env) PublisherRemove(c *gin.Context) {
	publisher := c.MustGet("publisher").(*models.Company)
	result := e.db.Where("name = ?", publisher.Name).Limit(1).Unscoped().Delete(&models.Company{}).GetErrors()

	if (models.HasError(result)) {
		e.ServerError(c)
		return
	}
	e.Json(c.Writer, map[string]interface{} {
		"ok": true,
	})
}

/**
 * List of publishers
 */
 func (e *Env) PublisherList(c *gin.Context) {
	limit := c.MustGet("limit").(int)
	offset := c.MustGet("offset").(int)
	limit = services.LimitRestriction(limit, "publishers")

	publishers := []models.Company{}
	result := e.db.Offset(offset).Limit(limit).Find(&publishers).GetErrors()
	if (models.HasError(result)) {
		e.ServerError(c)
		return
	}
	if (len(publishers) < 1) { // No publishers found
		e.Json(c.Writer, map[string]interface{} {
			"ok": true,
			"payload": []string{},
		})
		return
	}
	var ids []uint
	for _, publisher := range publishers {
		ids = append(ids, publisher.ID)
	}
	notifCounts, queryResult := models.GetDmInstance().CountNotificationsByPublishers(ids)
	if (models.HasError(queryResult)) {
		e.ServerError(c)
		return
	}
	var response []map[string]interface{}
	for _, publisher := range publishers {
		var ntfCount int
		var publisherId = int(publisher.ID)
		for _, notifCount := range *notifCounts {
			if (notifCount.PublisherId == publisherId) {
				ntfCount = notifCount.Total
			}
		}
		element := map[string]interface{}{ 
			"name": publisher.Name,
			"id": publisher.ID,
			"notifications_count": ntfCount,
		}
		response = append(response, element)
	}
	e.Json(c.Writer, map[string]interface{} {
		"ok": true,
		"payload": response,
	})
}

/**
 * Find total count of publishers
 */
func (e *Env) PublisherCount(c *gin.Context) {
	var count int
	e.db.Model(&models.Company{}).Count(&count)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": count,
	})
}

/**
 * Get publisher by ID
 */
func (e *Env) PublisherGetById(c *gin.Context) {
	publisher := c.MustGet("publisher").(*models.Company)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": gin.H{
			"id": publisher.ID,
			"name": publisher.Name,
		},
	})
}