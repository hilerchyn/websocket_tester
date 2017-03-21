// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)


func main() {

	// load config
	defaultConfig, _ := loadConfig()
	if defaultConfig == nil {
		log.Println("no configuration")
		os.Exit(1)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: defaultConfig.WSIP + ":" + defaultConfig.WSPort, Path: defaultConfig.WSPath}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Println(time.Now().Unix())
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	loginFlag := false

	for {
		select {
		case t := <-ticker.C:
			if loginFlag == false {
				log.Println(t.String())
				err := c.WriteMessage(websocket.TextMessage, []byte(defaultConfig.StrLogin))
				if err != nil {
					log.Println("write:", err)
					return
				}

				loginFlag = true
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
			c.Close()
			return
		}
	}
}
