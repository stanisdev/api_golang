package models

import(
	"app/services"
	"github.com/spf13/viper"
	"github.com/jinzhu/gorm"
	validator "github.com/asaskevich/govalidator"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"math/rand"
	"strings"
	"strconv"
	"fmt"
	"time"
)

type DbMethods struct {
	DB *gorm.DB
}

var instance *gorm.DB
var dmInstance *DbMethods

func DatabaseConnect() {
	params := viper.GetString("database.username") + ":" + viper.GetString("database.password") + "@/" + viper.GetString("database.dbname") + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", params)
  if err != nil {
    panic("Failed to connect database")
	}
	db.LogMode(true)
	fmt.Println("Database connected")
	instance = db
	dmInstance = &DbMethods{DB: instance}
}

func DatabaseMigrate() {
	instance.Set("gorm:table_options", "charset=utf8")
	instance.AutoMigrate(&User{}) // Create User table

	salt := services.GenerateRandomString(10)
	uniqKey := services.GenerateRandomString(10)
	hash, err := services.GetPasswordHash("pRCek5iFYm" + salt)
	if (err != nil) {
		panic(err)
	}
	user := User{}
	instance.Where(&User{Username: "mr.admin"}).Find(&user) // Create Admin record
	if (user.ID < 1) {
		admin := &User{Username: "mr.admin", Password: hash, Salt: salt, LastLogin: time.Now(), PasswordChanged: time.Now(), UniqUserKey: uniqKey}
		if err := instance.Create(admin); err != nil {
			fmt.Println("An error occurred while creating the \"Admin\" entry")
			fmt.Println(err)
		}
	}
	instance.AutoMigrate(&Notification{}) // Create Notification table
	instance.AutoMigrate(&Company{}) // Create Company table
	instance.Model(&Notification{}).AddForeignKey("company_id", "companies(id)", "SET NULL", "CASCADE")
}

func GetConnection() *gorm.DB {
	return instance
}

func GetDmInstance() *DbMethods {
	return dmInstance
}

func ValidateModel(modelInstance interface{}) bool {
	_, err := validator.ValidateStruct(modelInstance)
	return err == nil
}

func LoadFixtures() {
	db := GetConnection()
	publishersRaw := []string{ "Alibaba.com", "Flickr", "Instagram" }
	for _, publisher := range publishersRaw {
		db.Create(&Company{ Name: publisher })
	}
	publishers := []Company{}
	db.Find(&publishers)

	for a := 0; a <= 100; a++ {
		_a := strconv.Itoa(a)
		db.Create(&Notification{
			Message: "Message " + _a,
			Image: "Image" + _a,
			Header: "Header" + _a,
			Priority: 2,
			Expired: time.Now(),
			Button: "Button " + _a,
			Link: "http://link" + _a,
			CompanyID: publishers[rand.Intn(len(publishersRaw))].ID,
		})
	}
}

func HasError(errors []error) bool {
	if (len(errors) < 1) {
		return false
	}
	for _, err := range errors {
		errValue := strings.ToLower(err.Error())
		if (errValue == "record not found") {
			return false
		} else {
			fmt.Println(err.Error())
		}
	}
	return true
}