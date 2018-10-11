package services

import (
	"github.com/spf13/viper"
	"fmt"
	"os"
	"path/filepath"
	"path"
)

var dynamicConfig map[string]string
var curDir string

func DefineCurrentDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	curDir = dir
}

func ReadConfig() {
	DefineCurrentDir()
	confPath := os.Getenv("CONFIG_PATH")
	if len(confPath) < 1 || confPath == "*" {
		confPath = curDir
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(confPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func SetDynamicConfig() {
	uplPath := os.Getenv("UPLOADS_PATH")
	subPath := os.Getenv("SUB_PATH")

	if (len(uplPath) < 1) { // Or read from config file
		uplPath = viper.GetString("environment.uploads_dir")
	}
	if (len(subPath) < 1) {
		subPath = viper.GetString("environment.sub_path")
	}
	uplDir := os.Getenv("UPLOADS_DIR") // Path to "uploads" dir
	if (len(uplDir) < 1 || uplDir == "*") {
		uplDir = curDir
	}
	
	dynamicConfig = map[string]string{
		"UploadsDir": path.Join(uplDir, "uploads"),
		"UploadsPath": uplPath,
		"SubPath": subPath,
		"CurrentDir": curDir,
	}
}

func GetDynamicConfig() map[string]string {
	return dynamicConfig
}