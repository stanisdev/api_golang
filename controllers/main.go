package controllers

import (
	"app/models"
	"app/middlewares"
	"app/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"net/http"
	validator "github.com/asaskevich/govalidator"
	"encoding/json"
	"os"
	"fmt"
)

type Env struct {
	db *gorm.DB
	DBMethods *models.DbMethods
}

func (e *Env) Json(w http.ResponseWriter, data map[string]interface{}) {
	jsonOut, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonOut))
}

func (e *Env) ValidateStruct(structInstance interface{}) bool {
	_, err := validator.ValidateStruct(structInstance)
	return err == nil
}

func (e *Env) ServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"ok": false,
	})
}

func Start() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())

	subPath := services.GetDynamicConfig()["SubPath"]
	router.Static(subPath + "/uploads", services.GetDynamicConfig()["UploadsDir"])
	router.Use(middlewares.RequireAuthToken)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"ok": false,
			"message": "Page not found",
		})
	})
	env := &Env{
		db: models.GetConnection(),
		DBMethods: models.GetDmInstance(),
	}
	prefix := subPath + viper.GetString("environment.prefix")
	user := router.Group(prefix + "/user")
	{
		user.POST("/login", env.UserLogin)
		user.GET("/profile", env.UserProfile)
		user.POST("/change/password", env.UserChangePassword)
	}
	notification := router.Group(prefix + "/notification")
	{
		notification.GET("/list", middlewares.LimitOffset, env.NotificationList)
		notification.POST("/create", middlewares.ValidateNotification, env.NotificationCreate)
		notification.GET("/delete/:id", middlewares.UrlIdCorrectness, middlewares.FindNotificationById, env.NotificationRemove)
		notification.GET("/get/:id", middlewares.UrlIdCorrectness, middlewares.FindNotificationById, env.NotificationGetById)
		notification.POST("/edit/:id", middlewares.UrlIdCorrectness, middlewares.FindNotificationById, middlewares.ValidateNotification, env.NotificationUpdate)
		notification.GET("/count", env.NotificationCount)
	}
	image := router.Group(prefix + "/image")
	{
		image.POST("/upload", env.ImageUpload)
	}
	publisher := router.Group(prefix + "/publisher")
	{
		publisher.POST("/create", env.PublisherCreate)
		publisher.POST("/edit/:id", middlewares.FindPublisherById, env.PublisherUpdate)
		publisher.GET("/delete/:id", middlewares.FindPublisherById, env.PublisherRemove)
		publisher.GET("/list", middlewares.LimitOffset, env.PublisherList)
		publisher.GET("/count", env.PublisherCount)
		publisher.GET("/get/:id", middlewares.FindPublisherById, env.PublisherGetById)
		publisher.GET("/plain_list", env.PublisherPlainList)
	}
	router.GET(subPath + "/notifications", env.NotificationPublic)

	_port := os.Getenv("PORT")
	port := viper.GetString("environment.port")
	if (len(_port) > 0) {
		port = _port
	}
	fmt.Println("App are listening port " + port)
	
	router.Run(":" + port)
}