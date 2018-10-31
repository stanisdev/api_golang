package controllers

import (
	_ "app/models"
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
	e.Json(c.Writer, map[string]interface{} {
		"ok": true,
	})
}

/**
 * Edit publisher
 */
func (e *Env) PublisherUpdate(c *gin.Context) {

}