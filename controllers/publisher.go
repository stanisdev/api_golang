package controllers

import (
	"app/models"
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
	e.Json(c.Writer, map[string]interface{} {
		"ok": true,
	})
}