package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var AppPath string
var AppConfigC AppConfig

func InitConfig() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	AppPath = path
	configFilePath := path + "/config.yaml"

	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal("open config file err:", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfigC)
	if err != nil {
		log.Fatal("decode config err:", err)
	}

	log.Print("Config | Load Config success~")
}
