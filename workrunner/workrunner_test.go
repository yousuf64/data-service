package workrunner

import (
	"sync"
	"testing"
	"time"
)

func TestWorkRunner_Run(t *testing.T) {
	wr := New[string, string]()

	var wg sync.WaitGroup
	iters := 10
	wg.Add(iters)

	for i := 0; i < iters; i++ {
		go func() {
			response := <-wr.Observe("911_365", func() string {
				time.Sleep(time.Millisecond * 1)
				return "111"
			})

			t.Log(response)
			wg.Done()
		}()
	}

	wg.Wait()
}
