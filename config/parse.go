package config

import (
	"os"
	"errors"
	"encoding/json"
)

type Config struct {
	WSIP     string `json:"ws_ip"`
	WSPort   string `json:"ws_port"`
	WSPath   string `json:"ws_path"`
	StrLogin string `json:"str_login"`
	StrSay   string `json:"str_say"`
}


func NewConfig(configFile string) (*Config, error) {

	file, err := os.Open(configFile)
	if err != nil {
		return nil, errors.New("open file failed!")
	}

	config := Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, errors.New("decode config failed")
	}

	return &config, nil

}