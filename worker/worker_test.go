package worker

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var counter atomic.Int32

func TestWorker(t *testing.T) {
	vals := []string{
		"mon",
		"tue",
		"wed",
		"thu",
		"fri",
		"sat",
		"sun",
	}

	wkr := New[string](func() string {
		log.Println("executing")
		time.Sleep(time.Millisecond * 2)
		defer counter.Add(1)
		return vals[counter.Load()]
	})

	var wg sync.WaitGroup
	wg.Add(100_000)

	for i := 0; i < 100_000; i++ {
		wkr.Subscribe(func(s string) {
			wg.Done()
		})
		//time.Sleep(time.Millisecond * 320)
	}

	wg.Wait()
}
