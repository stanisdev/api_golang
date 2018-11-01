package middlewares

import (
	"app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "fmt"
)

func FindPublisherById(c *gin.Context) {
	id := c.Param("id")
	publisher := &models.Company{}
	models.GetConnection().Where("id = ?", id).First(publisher)

	if (publisher.ID < 1) {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
		})
		c.Abort()
		return
	}
	c.Set("publisher", publisher)
	c.Next()
}