package services

import (
	"github.com/spf13/viper"
	"fmt"
	"os"
	"path/filepath"
	"path"
)

var dynamicConfig map[string]string

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./src/app")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}	
}

func SetDynamicConfig() {
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if (rootDir[len(rootDir)-4:] != "/bin") {
		panic("Wrong env structure")
	}
	rootDir = rootDir[0:len(rootDir)-4]
	appDir := path.Join(rootDir, "src", "app")
	uplPath := os.Getenv("UPLOADS_PATH")
	subPath := os.Getenv("SUB_PATH")

	if (len(uplPath) < 1) { // Or read from config file
		uplPath = viper.GetString("environment.uploads_dir")
	}
	if (len(subPath) < 1) {
		subPath = viper.GetString("environment.sub_path")
	}

	dynamicConfig = map[string]string{
		"RootDir": rootDir,
		"AppDir": appDir,
		"UploadsDir": path.Join(appDir, "uploads"),
		"UploadsPath": uplPath,
		"SubPath": subPath,
	}
}

func GetDynamicConfig() map[string]string {
	return dynamicConfig
}