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
	

	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(s.defaultConfig.StrLogin, rand.Int())))
	if err != nil {
		// release one wait chan
		<-waitChan
		log.Println("write:", err)
		return
	}
	s.TotalConn++

	// release one wait chan
	<-waitChan
	
	ticker := time.NewTicker(time.Duration(s.defaultConfig.ExecSecond) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			return
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
				log.Printf("recv[%d]: %s", workerId, message)
			}
		}
	}

}
