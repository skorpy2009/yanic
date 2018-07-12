package all

import (
	"sync"
	"time"

	"chaos.expert/FreifunkBremen/yanic/output"
	"chaos.expert/FreifunkBremen/yanic/runtime"
)

var quit chan struct{}
var wg = sync.WaitGroup{}
var outputA output.Output

func Start(nodes *runtime.Nodes, config runtime.NodesConfig) (err error) {
	outputA, err = Register(config.Output)
	if err != nil {
		return
	}
	quit = make(chan struct{})
	wg.Add(1)
	go saveWorker(nodes, config.SaveInterval.Duration)
	return
}

func Close() {
	close(quit)
	wg.Wait()
	quit = nil
}

// save periodically to output
func saveWorker(nodes *runtime.Nodes, saveInterval time.Duration) {
	ticker := time.NewTicker(saveInterval)
	for {
		select {
		case <-ticker.C:
			outputA.Save(nodes)
		case <-quit:
			ticker.Stop()
			wg.Done()
			return
		}
	}
}
