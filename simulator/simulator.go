package simulator

import (
	"net/url"


	"chat_tester/config"
	"sync"
)

type Simulator struct {
	defaultConfig *config.Config
	Url url.URL
	Count int
	wg sync.WaitGroup
	worker map[int]chan int
}

func NewSimulator(defaultConfig *config.Config) (*Simulator, error) {

	// url
	u := url.URL{Scheme: defaultConfig.WSScheme, Host: defaultConfig.WSIP + ":" + defaultConfig.WSPort, Path: defaultConfig.WSPath}

	return &Simulator{defaultConfig:defaultConfig, Url:u, Count:defaultConfig.SimulatorCount}, nil
}


func (s *Simulator) Run(){

	// init worker map
	s.worker = make(map[int] chan int)

	// start worker
	for count := 0; count < s.Count ; count++ {
		s.wg.Add(1)
		go s.connect(count)
	}

	s.wg.Wait()
}
