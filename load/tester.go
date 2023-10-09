package load

import (
	"load-test/net"
	"sync"
	"time"
)

type TesterConfig struct {
	NumberUsers int
	Rampup      int
	Duration    int
	Target      string
	Path        string
}

type Tester struct {
	config TesterConfig
	client net.Client
}

func NewTester(config TesterConfig, client net.Client) *Tester {
	return &Tester{config, client}
}

func (tester *Tester) Run() error {
	rampupTime := time.Duration(tester.config.Rampup) * time.Second
	rampupStep := rampupTime / time.Duration(tester.config.NumberUsers)

	errors := make(chan error)
	running := true

	go func() {
		select {
		case <-time.After(time.Duration(tester.config.Duration) * time.Second):
			running = false
			close(errors)
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < tester.config.NumberUsers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for running {
				beforeRequest := time.Now()
				if err := tester.client.Send(tester.config.Target, tester.config.Path); err != nil {
					errors <- err
				}

				dt := time.Since(beforeRequest)
				waitTime := 1*time.Second - dt
				time.Sleep(waitTime)
			}
		}()

		time.Sleep(rampupStep)
	}

	for err := range errors {
		return err
	}

	wg.Wait()
	return nil
}
