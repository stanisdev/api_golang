package controllers

import (
	"app/models"
	"app/services"
	structs "app/structures"
	"github.com/gin-gonic/gin"
	"time"
	_ "fmt"
)

func (e *Env) UserLogin(c *gin.Context) {
	var login structs.UserLogin
	c.BindJSON(&login)
		
	if (len(login.Username) < 1 || len(login.Password) < 1) {
		services.WrongPostData(c)
		return
	}
	user, exists := e.DBMethods.UserFindOne(&models.User{Username: login.Username})
	wrongCred := func () {
		c.JSON(200, gin.H{
			"ok": false,
			"errors": gin.H{
				"pass": "Wrong password/username",
			},
		})
	}
	if (!exists) {
		wrongCred()
		return
	}
	isValid := services.CheckPasswordHash(login.Password + user.Salt, user.Password)
	if (!isValid) {
		wrongCred()
		return
	}
	tokenString, err := services.CryptJWT(user.ID, user.UniqUserKey)
	if (err != nil) {
		c.JSON(200, gin.H{
			"ok": false,
			"message": "Can not log in",
		})
		return
	}
	user.LastLogin = time.Now()
	e.db.Save(&user)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": gin.H{
			"id_token": tokenString,
		},
	})
}

func (e *Env) UserProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.JSON(200, gin.H{
		"ok": true,
		"payload": gin.H{
			"username": user.Username,
			"created_ad": user.CreatedAt.Unix(),
			"last_login": user.LastLogin.Unix(),
			"password_changed": user.PasswordChanged.Unix(),
		},
	})
}

func (e *Env) UserChangePassword(c *gin.Context) {
	var chp structs.UserChangePassword
	c.BindJSON(&chp)

	if (len(chp.OldPassword) < 1 || len(chp.NewPassword) < 1 || chp.OldPassword == chp.NewPassword) {
		services.WrongPostData(c)
		return
	}
	user := c.MustGet("user").(models.User)
	isValid := services.CheckPasswordHash(chp.OldPassword + user.Salt, user.Password)
	if (!isValid) {
		c.JSON(200, gin.H{
			"ok": false,
			"errors": gin.H{
				"old_password": "Wrong password",
			},
		})
		return
	}
	salt := services.GenerateRandomString(10)
	hash, err := services.GetPasswordHash(chp.NewPassword + salt)
	if (err != nil) {
		services.ServerError(err, c)
		return
	}
	user.Salt = salt
	user.Password = hash
	user.PasswordChanged = time.Now()
	e.db.Save(&user)
	
	c.JSON(200, gin.H{
		"ok": true,
	})
}