package main

import (
	"flag"
	"github.com/hilerchyn/websocket_tester/config"
	"log"
	"os"
)

// flag
var cfg = flag.String("cfg", "config.json", "access path")

func loadConfig() (*config.Config, error) {
	flag.Parse()
	log.SetFlags(0)

	defaultConfig, err := config.NewConfig(*cfg)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	return defaultConfig, nil
}
