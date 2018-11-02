package middlewares

import (
	"app/services"
	"app/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	_ "fmt"
)

type Notification struct {
	Message string `json:"message"`
	Image string `json:"image"`
	Header string `json:"header"`
	Priority string `json:"priority"`
	Expired string `json:"expired"`
	Button string `json:"button"`
	Link string `json:"link"`
	PublisherId string `json:"publisher_id"`
}

func ValidateNotification(c *gin.Context) {
	var ntf Notification
	var pubId int

	c.BindJSON(&ntf)

	if _pubId, err := strconv.Atoi(ntf.PublisherId); err != nil {
		services.WrongPostData(c)
		return
	} else {
		pubId = _pubId
	}

	publ := &models.Company{}
	models.GetConnection().Where("id = ?", pubId).First(&publ) // Find publisher by ID
	if (publ.ID < 1) {
		services.WrongPostData(c)
		return
	}
	exp, err0 := time.Parse("2006/01/02", ntf.Expired) // Parse expired data
	if err0 != nil {
		services.WrongPostData(c)
		return
	}
	prior, err1 := strconv.ParseUint(ntf.Priority, 10, 32)
	if err1 != nil {
		services.WrongPostData(c)
		return
	}
	ntfInstance := &models.Notification { // Creating Notification
		Message: ntf.Message,
		Image: ntf.Image,
		Header: ntf.Header,
		Priority: uint(prior),
		Expired: exp,
		Button: ntf.Button,
		Link: ntf.Link,
		CompanyID: uint(pubId),
	}
	
	if (!models.ValidateModel(ntfInstance)) {
		services.WrongPostData(c)
		return
	}
	c.Set("notificationBlank", ntfInstance)
	c.Next()
}