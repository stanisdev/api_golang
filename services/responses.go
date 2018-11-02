package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func WrongPostData(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"ok": false,
		"message": "Wrong post data",
	})
	c.Abort()
}

func WrongUrlParams(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"ok": false,
		"message": "Wrong URL params",
	})
}

func ServerError(err error, c *gin.Context) {
	fmt.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"ok": false,
		"message": "Server Error",
	})
}