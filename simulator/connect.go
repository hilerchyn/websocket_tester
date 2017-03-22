package simulator

import (
	"log"
	"time"

	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"strings"
	"sync"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s *Simulator) connect(workerId int) {

	log.Printf("worker %d connecting to %s", workerId, s.Url.String())
	//c := func() *websocket.Conn {
	//	for {
	c, _, err := websocket.DefaultDialer.Dial(s.Url.String(), nil)
	if err != nil {
		s.wg.Done()
		log.Println("dial:", err)
		return
	}

	//		return c
	//	}
	//}()
	defer c.Close()

	// set chan value
	s.worker[workerId] = make(chan int)

	// read message
	go s.sync(workerId, c)

	lock := new(sync.Mutex)

	loginFlag := false
	tickerLogin := time.NewTicker(time.Second + time.Duration(rand.Intn(s.defaultConfig.SimulatorStartIn)+1))
	defer tickerLogin.Stop()
	ticker := time.NewTicker(time.Duration(s.defaultConfig.ExecSecond) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			//log.Println(t.String())
			s.worker[workerId] <- workerId
			return
		case <-tickerLogin.C:
			if loginFlag == false {
				if c == nil {
					continue
				}
				//log.Println(tl.String())
				lock.Lock()
				err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(s.defaultConfig.StrLogin, rand.Int())))
				lock.Unlock()
				if err != nil {
					s.worker[workerId] <- workerId
					log.Println("write:", err)
					return
				}
				s.TotalConn++
				loginFlag = true
			}
		default:

		}
	}

}

func (s *Simulator) sync(workerId int, conn *websocket.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	lock := new(sync.Mutex)

	for {
		select {
		case done := <-s.worker[workerId]:
			log.Println("Worker:", done, " Done!")
			close(s.worker[workerId])
			return
		default:
			if conn == nil {
				continue
			}

			lock.Lock()
			_, message, err := conn.ReadMessage()
			lock.Unlock()
			if err != nil {
				log.Println("read[", workerId, "]:", err)
				return
			}
			log.Print(time.Now().String())
			log.Printf("recv[%d]: %s", workerId, message)

			// pong
			if strings.Compare(string(message), s.defaultConfig.StrPing) == 0 {
				lock.Lock()
				err := conn.WriteMessage(websocket.TextMessage, []byte(s.defaultConfig.StrPong))
				lock.Unlock()
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}

	}
}
