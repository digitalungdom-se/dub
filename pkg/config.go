package pkg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Prefix  []string `json:"prefix"`
	GuildID string   `json:"guildID"`
}

func LoadConfig() Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("error opening config file")
	}

	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	return config
}
