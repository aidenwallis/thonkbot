package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	TwitchUsername string `json:"twitchUsername"`
	TwitchPassword string `json:"twitchPassword"`
	MysqlDSN       string `json:"mysqlDsn"`
	Webhost        string `json:"webhost"`
}

var Cfg *Config

func LoadConfig() {
	bs, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	Cfg = &Config{}

	err = json.Unmarshal(bs, Cfg)
	if err != nil {
		panic(err)
	}
}
