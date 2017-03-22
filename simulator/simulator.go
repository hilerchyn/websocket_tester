package simulator

import (
	"log"
	"net/url"

	"github.com/hilerchyn/websocket_tester/config"
	"os"
	"sync"
)

type Simulator struct {
	defaultConfig *config.Config
	Url           url.URL
	Count         int
	wg            sync.WaitGroup
	worker        map[int]chan int
	TotalConn     int
}

func NewSimulator(defaultConfig *config.Config) (*Simulator, error) {

	if defaultConfig.SimulatorStartIn >= defaultConfig.ExecSecond {
		log.Println("simulator_start_in should greater than exec_second")
		os.Exit(1)
	}

	// url
	u := url.URL{Scheme: defaultConfig.WSScheme, Host: defaultConfig.WSIP + ":" + defaultConfig.WSPort, Path: defaultConfig.WSPath}

	return &Simulator{defaultConfig: defaultConfig, Url: u, Count: defaultConfig.SimulatorCount}, nil
}

func (s *Simulator) Run() {

	// init worker map
	s.worker = make(map[int]chan int)

	s.TotalConn = 0

	// start worker
	for count := 0; count < s.Count; count++ {
		s.wg.Add(1)
		go s.connect(count)
	}

	s.wg.Wait()

	log.Println("TotalConnections:", s.TotalConn)
}
