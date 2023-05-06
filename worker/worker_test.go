package worker

import (
	"github.com/hashicorp/logutils"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func init() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "DEBUG",
		Writer:   os.Stdout,
	}
	log.SetOutput(filter)
}

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
		"foo",
		"bar",
		"abc",
		"qqq",
		"zzz",
		"vvv",
		"bbb",
		"mmm",
	}

	wkr := New[string](func() string {
		log.Println("executing")
		time.Sleep(time.Nanosecond * 500)
		defer counter.Add(1)
		return vals[counter.Load()]
	})

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		wkr.Subscribe(func(s string) {
			wg.Done()
		})
		//time.Sleep(time.Millisecond * 320)
	}

	wg.Wait()
}
