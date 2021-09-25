package worker

import (
	"errors"
	"time"
)

// TODO: move to config.Config
const maxGoroutines = 10

type Func func(arg interface{}) error

func Worker(arg interface{}, fn func(arg interface{}) error) error {
	// TODO: Stop worker when error is returned

	// queue control to number of Goroutines.
	queue := make(chan struct{}, maxGoroutines)
	defer close(queue)

	errChan := make(chan struct{})
	defer close(errChan)

	//var wg sync.WaitGroup

	for {
		queue <- struct{}{}

		select {
		case <-errChan:
			return errors.New("something happened")
		default:
			go func() {
				time.Sleep(200 * time.Millisecond)
				//defer wg.Done()
				defer func() { <-queue }()

				err := fn(arg)
				if err != nil {
					errChan <- struct{}{}
				}
			}()
		}
	}
}
