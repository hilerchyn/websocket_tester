package simulator

import (
	"log"
	"time"

	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"strings"
	//	"sync"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s *Simulator) connect(workerId int, waitChan chan int) {

	defer s.wg.Done()

	log.Printf("worker %d connecting to %s", workerId, s.Url.String())

	c, _, err := websocket.DefaultDialer.Dial(s.Url.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		// release one wait chan
		<-waitChan
		return
	}
	defer c.Close()
	

	// generate user ID
	userID := rand.Int()
	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(s.defaultConfig.StrLogin, userID, userID)))
	if err != nil {
		// release one wait chan
		<-waitChan
		log.Println("write:", err)
		return
	}
	s.TotalConn++

	// release one wait chan
	<-waitChan

	var sayTicker *time.Ticker
	var sayString string
	if s.defaultConfig.SayInterval == 0 {
		sayTicker = time.NewTicker(time.Duration(s.defaultConfig.ExecSecond+10) * time.Second)
	} else {

		sayString = fmt.Sprintf(s.defaultConfig.StrSay, userID)
		sayTicker = time.NewTicker(time.Duration(s.defaultConfig.SayInterval + rand.Intn(s.defaultConfig.SayInterval)) * time.Second)
	}

	ticker := time.NewTicker(time.Duration(s.defaultConfig.ExecSecond) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			return
		case <-sayTicker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(sayString))
			if err != nil {
				log.Println("write:", err)
				return
			}
		default:

			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read[", workerId, "]:", err)
				return
			}


			if strings.Compare(string(message), s.defaultConfig.StrPing) == 0 {
				err := c.WriteMessage(websocket.TextMessage, []byte(s.defaultConfig.StrPong))
				if err != nil {
					log.Println("write:", err)
					return
				}
			} else {
				if s.defaultConfig.OutputRecv {
					log.Printf("recv[%d]: %s", workerId, message)
				}
			}
		}
	}

}
