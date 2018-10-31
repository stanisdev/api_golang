package controllers

import (
	"app/services"
	"github.com/gin-gonic/gin"
	"os"
	"io"
	"path"
	_ "fmt"
	"errors"
)

/**
 * Image uploading
 */
func (e *Env) ImageUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if (err != nil) {
		services.WrongPostData(c)
		return
	}
	fileName := header.Filename
	ext := path.Ext(fileName)[1:]
	allowedExt := []string{"jpg", "jpeg", "png", "gif"}
	correctExt := false
	for _, _ext := range allowedExt {
		if (_ext == ext) {
			correctExt = true
			break
		}
	}
	if (!correctExt) {
		services.ServerError(errors.New("Not image file"), c)
		return
	}

	fileName = services.GenerateRandomString(16) + "." + ext
	dest := path.Join(services.GetDynamicConfig()["UploadsDir"], fileName)

	out, err := os.Create(dest)
	if err != nil {
		services.ServerError(err, c)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		services.ServerError(err, c)
		return
	}

	c.JSON(200, gin.H{
		"ok": true,
		"payload": fileName,
	})
}