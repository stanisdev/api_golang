package controllers

import (
	"app/models"
	"app/services"
	"github.com/gin-gonic/gin"
	structs "app/structures"
	"net/http"
	"fmt"
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
	publisher.Name = publ.Name
	result := e.db.Save(&publisher).GetErrors()
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
	publishers := []models.Company{}
	result := e.db.Find(&publishers).GetErrors()
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
	fmt.Println(ids)

	e.Json(c.Writer, map[string]interface{} {
		"ok": true,
	})
}