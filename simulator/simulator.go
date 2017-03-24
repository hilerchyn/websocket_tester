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
	//worker        map[int]chan int
	TotalConn     int
}

func NewSimulator(defaultConfig *config.Config) (*Simulator, error) {

	// validate parameters
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
	//s.worker = make(map[int]chan int)

	// statistic how many connections success
	s.TotalConn = 0

	// set how many workers to start the connection at the same time
	waitChan := make(chan int, s.defaultConfig.WorkerCount)

	// start worker
	for count := 0; count < s.Count; count++ {
		s.wg.Add(1)
		waitChan <-count
		go s.connect(count, waitChan)
	}

	// waitting for all connections end
	s.wg.Wait()

	// out put how many connections started
	log.Println("TotalConnections:", s.TotalConn)
}
