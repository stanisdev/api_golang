package main

import (
	"app/controllers"
	"app/services"
	"app/models"
	"math/rand"
	"time"
	"os"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	services.ReadConfig()
	services.SetDynamicConfig()
	models.DatabaseConnect()	
	models.DatabaseMigrate()

	if (len(os.Getenv("LOAD_FIXTURES")) > 0) {
		models.LoadFixtures()
	} else {
		controllers.Start()
	}
}