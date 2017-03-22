// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"log"
	"os"

	"github.com/hilerchyn/websocket_tester/simulator"
)

func main() {

	// load config
	defaultConfig, _ := loadConfig()
	if defaultConfig == nil {
		log.Println("no configuration")
		os.Exit(1)
	}

	// simulator
	simulator, err := simulator.NewSimulator(defaultConfig)
	if err != nil {
		log.Println("create simulator failed")
		os.Exit(1)
	}

	// start
	simulator.Run()

}
