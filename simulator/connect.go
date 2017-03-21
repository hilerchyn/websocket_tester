package simulator

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"strings"
)

func (s *Simulator) test(workerId int){

	log.Printf("worker %d connecting to %s", workerId, s.Url.String())
	c, _, err := websocket.DefaultDialer.Dial(s.Url.String(), nil)
	if err != nil {
		s.wg.Done()
		log.Fatal("dial:", err)
	}
	//defer c.Close()

	// set chan value
	s.worker[workerId] = make(chan int)

	// read message
	go s.sync(workerId, c)

	loginFlag := false
	tickerLogin := time.NewTicker(time.Second)
	defer tickerLogin.Stop()
	ticker := time.NewTicker(time.Duration(s.defaultConfig.ExecSecond) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			log.Println(t.String())
			s.worker[workerId] <- workerId
			return
		case tl := <-tickerLogin.C:
			if loginFlag == false {
				log.Println(tl.String())
				err := c.WriteMessage(websocket.TextMessage, []byte(s.defaultConfig.StrLogin))
				if err != nil {
					log.Println("write:", err)
					return
				}

				loginFlag = true
			}
		default:

		}
	}

}

func (s *Simulator) sync(workerId int, conn *websocket.Conn){
	defer s.wg.Done()
	defer conn.Close()

	for {
		select {
		case done := <- s.worker[workerId]:
			log.Println("Worker:", done, " Done!")
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Println(time.Now().Unix())
			log.Printf("recv: %s", message)

			// pong
			if strings.Compare(string(message), s.defaultConfig.StrPing) == 0 {
				err := conn.WriteMessage(websocket.TextMessage, []byte(s.defaultConfig.StrPong))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}

	}
}
