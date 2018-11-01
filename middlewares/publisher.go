package middlewares

import (
	"app/models"
	"app/services"
	"github.com/gin-gonic/gin"
	"strconv"
	_ "fmt"
)

func FindPublisherById(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.ParseInt(id, 10, 32); err != nil {
		services.WrongUrlParams(c)
		c.Abort()
		return
	}
	publisher := &models.Company{}
	models.GetConnection().Where("id = ?", id).First(publisher)

	if (publisher.ID < 1) {
		services.WrongUrlParams(c)
		c.Abort()
		return
	}
	c.Set("publisher", publisher)
	c.Next()
}