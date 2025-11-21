package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"simple-file-server/global"
	"simple-file-server/lib/defs"
)

func Init() {
	config := getConfig()
	global.CONFIG = config
	log.WithFields(log.Fields{
		"port": config.Port,
	}).Info("config")
}

func getConfig() defs.Config {
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		file, err := os.Create("config.json")
		if err != nil {
			log.Fatal(err)
		}
		file.Chmod(0700)
		defer file.Close()
		config := defs.Config{
			Debug:    false,
			ApiToken: "xxx",
			Port:     60088,
			DataDir:  "./data",
			TempDir:  "./temp",
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		err = encoder.Encode(config)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Default config file created at config.json")
	}
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := defs.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	if config.TempDir == "" {
		config.TempDir = "./temp"
	}
	if config.DataDir == "" {
		config.DataDir = "./data"
	}
	return config
}
